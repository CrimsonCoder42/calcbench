// calculator_net_run.mjs: submit one input set to one calculator page and print the result.
//
// Usage: node scripts/calculator_net_run.mjs <calculator URL> '<JSON of field name -> value>' [form name] [button label]
// form name "" means the whole page; button label picks one of several buttons by its text (e.g. "Multiply").
// Field values: text and select take a string (a select matches by option value, then by option label);
// a radio takes the option value; a checkbox takes true/false. Fields not named keep the page default.
// A field is named by its name attribute, or by its id when it has none. One page load and one submit (the
// submit is a second navigation; the waits between loads are the agent's job). A page with no Calculate button
// is treated as live: typing events are fired and the page read. Prints JSON: {url, status, submit
// (navigated | in-place | live), inputs, sent, missing, result: {text, tables}, fieldsAfter (every form text
// box and select after the submit, where solvers write their answers)}. Only the result block
// is read: elements whose class names contain "result" (or the #printit block), their text and their tables.
// Exit 3 on 403/429 or a block page, 4 when no calculate button or no result block is found.

import { chromium } from '@playwright/test';

const [url, json, formName, buttonLabel] = process.argv.slice(2);   // formName "" = whole page
if (!url || !json) {
  console.error("usage: calculator_net_run.mjs <calculator URL> '<JSON field->value>' [form name]");
  process.exit(2);
}
const inputs = JSON.parse(json);

const blocked = (status, title) =>
  status === 403 || status === 429 || /access denied|attention required|just a moment|captcha/.test(title.toLowerCase());

const browser = await chromium.launch({ headless: true });
let code = 0;
try {
  const page = await browser.newPage();
  const resp = await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });
  // Some calculators draw their form from script after the first paint.
  await page.waitForLoadState("load", { timeout: 15000 }).catch(() => {});
  await page.waitForTimeout(1500);
  const status = resp ? resp.status() : 0;
  if (blocked(status, await page.title())) {
    console.log(JSON.stringify({ url, status, blocked: true }));
    code = 3;
  } else {
    const scope = formName ? page.locator(`form[name="${formName}"]`).first() : page;
    const formSel = "";   // fields are found anywhere in the scope, in or out of a form
    const missing = [];
    // Selects, radios and checkboxes first, text last: a unit select's change handler may rewrite the
    // number beside it, so the number is typed after the unit is set.
    // A field is found by its name, or by its id when it has no name (some date pickers).
    const field = name => scope.locator(`${formSel}[name="${name}"], ${formSel}[id="${name}"]`);
    const kind = async name => {
      const l = field(name);
      if (await l.count() === 0) return 0;
      return l.first().evaluate(el => (el.tagName === "SELECT" || el.type === "radio" || el.type === "checkbox") ? 0 : 1);
    };
    const order = [];
    for (const [name, value] of Object.entries(inputs)) order.push([await kind(name), name, value]);
    order.sort((a, b) => a[0] - b[0]);
    for (const [, name, value] of order) {
      const loc = field(name);
      const n = await loc.count();
      if (n === 0) { missing.push(name); continue; }
      const first = loc.first();
      const tag = await first.evaluate(el => el.tagName);
      const type = await first.evaluate(el => (el.type || '').toLowerCase());
      if (tag === 'SELECT') {
        const v = String(value);
        const ok = await first.evaluate((el, v) => {
          const o = [...el.options].find(o => o.value === v) || [...el.options].find(o => o.textContent.trim() === v);
          if (!o) return false;
          if (el.value === o.value) return true;
          el.value = o.value;
          el.dispatchEvent(new Event('change', { bubbles: true }));
          return true;
        }, v);
        if (!ok) missing.push(`${name}=${v}`);
      } else if (type === 'radio') {
        const r = scope.locator(`${formSel}input[type=radio][name="${name}"][value="${value}"]`);
        if (await r.count() === 0) missing.push(`${name}=${value}`);
        else await r.first().evaluate(el => { if (!el.checked) el.click(); });
      } else if (type === 'checkbox') {
        // Styled checkboxes sit outside the viewport; a DOM click still runs the page's onclick.
        await first.evaluate((el, v) => { if (el.checked !== v) el.click(); }, Boolean(value));
      } else {
        if (!(await first.isVisible())) {
          await first.evaluate((el, v) => { el.value = v; }, String(value));
        } else {
          await first.fill(String(value));
          // Some pages copy a typed box into the submitted field only on change or blur.
          await first.evaluate(e => { for (const t of ["keyup", "change"]) e.dispatchEvent(new Event(t, { bubbles: true })); e.dispatchEvent(new Event("blur")); });
        }
      }
    }
    // Tried in priority order (a comma list would return the first in page order, e.g. the site search button).
    // With no form named, the button is looked for first in the form that holds the first field set (a page
    // with several calculators has one Calculate button per form), then on the whole page.
    const scopes = [];
    if (!formName && Object.keys(inputs).length) {
      const f = field(Object.keys(inputs)[0]);
      if (await f.count()) {
        const idx = await f.first().evaluate(e => { const fm = e.closest('form'); return fm ? [...document.forms].indexOf(fm) : -1; });
        if (idx >= 0) scopes.push(page.locator('form').nth(idx));
      }
    }
    scopes.push(scope);
    let button = null;
    const wanted = buttonLabel ? [`input[type=submit][value*="${buttonLabel}"]`, `input[type=button][value*="${buttonLabel}"]`,
      `button:has-text("${buttonLabel}")`] : [];
    for (const sc of scopes) {
      for (const sel of [...wanted, 'input[type=submit][value*="alculate"]', 'input[type=button][value*="alculate"]',
        'input[type=image][alt*="alculate"]', 'button:has-text("Calculate")', 'input[type=submit]', 'input[type=image]']) {
        const l = sc.locator(sel);
        const n = await l.count();
        for (let i = 0; i < n && !button; i++) if (await l.nth(i).isVisible()) button = l.nth(i);
        if (button) break;
      }
      // A labelled control that is not a button element (a link or a styled cell): match its exact text.
      if (!button && buttonLabel) {
        const t = sc.getByText(buttonLabel, { exact: true });
        const n = await t.count();
        for (let i = 0; i < n && !button; i++) if (await t.nth(i).isVisible()) button = t.nth(i);
      }
      if (button) break;
    }
    const live = button === null;
    {
      // What the form will send, read back after filling (catches a field the page reformatted or ignored).
      const sent = {};
      for (const name of Object.keys(inputs)) {
        const el = field(name);
        if (await el.count() === 0) continue;
        sent[name] = await el.first().evaluate(e => e.type === "checkbox" ? e.checked
          : e.type === "radio" ? ((e.form || document).querySelector(`input[name="${e.name}"]:checked`) || {}).value : e.value);
      }
      const before = page.url();
      if (live) {
        // No Calculate button: the page computes as you type. Fire the typing events on every field set.
        for (const name of Object.keys(inputs)) {
          const el = field(name);
          if (await el.count() === 0) continue;
          await el.first().evaluate(e => { for (const t of ["input", "keyup", "change"]) e.dispatchEvent(new Event(t, { bubbles: true })); });
        }
        await page.waitForTimeout(1500);
      } else {
        await button.click();
        // A server-side calculator navigates; a client-side one rewrites the page in place.
        await page.waitForURL(u => u.href !== before, { timeout: 10000 }).catch(() => {});
        await page.waitForLoadState("domcontentloaded");
        await page.waitForTimeout(1000);
      }
      const st2 = live ? "live" : page.url() !== before ? "navigated" : "in-place";
      // Calculators that solve for the missing quantities write them back into the form's own boxes.
      const fieldsAfter = await page.evaluate(() => {
        const out = {};
        for (const e of document.querySelectorAll("input[type=text], input[type=number], input:not([type]), select")) {
          const k = e.name || e.id;
          if (k && !(k in out)) out[k] = e.value;
        }
        return out;
      });
      if (blocked(0, await page.title())) {
        console.log(JSON.stringify({ url, status, blocked: true, after: 'submit' }));
        code = 3;
      } else {
        const result = await page.evaluate(() => {
          const clean = s => (s || '').replace(/[ \t]+/g, ' ').replace(/\n\s*\n+/g, '\n').trim();
          const visible = el => !!(el.offsetWidth || el.offsetHeight || el.getClientRects().length);
          let found = [...document.querySelectorAll('[class*="result" i], [id*="result" i], #printit')]
            .filter(visible);
          found = found.filter(r => !found.some(o => o !== r && o.contains(r)));   // outermost only
          // A result heading is followed by its figures and tables as siblings: take those too, up to the
          // next plain heading or anything holding a form field.
          const roots = [];
          for (const r of found) {
            roots.push(r);
            if (!/^H[1-4]$/.test(r.tagName)) continue;
            let s = r.nextElementSibling;
            for (let i = 0; s && i < 40; i++, s = s.nextElementSibling) {
              if (/^H[1-3]$/.test(s.tagName) && !/result/i.test(s.className)) break;
              if (s.matches("form") || s.querySelector("input, select, form")) break;
              if (visible(s) && !found.includes(s)) roots.push(s);
            }
          }
          const text = roots.map(r => clean(r.innerText)).filter(Boolean);
          const tables = [];
          for (const r of roots) {
            for (const t of r.querySelectorAll('table')) {
              if (t.querySelector('table') || !visible(t)) continue;               // leaf tables only
              const rows = [...t.rows].map(row => [...row.cells].map(c => clean(c.innerText)));
              if (rows.length) tables.push(rows.slice(0, 400));
            }
          }
          return { text, tables };
        });
        const empty = result.text.length === 0;
        console.log(JSON.stringify({ url, status, submit: st2, inputs, sent, missing, result, fieldsAfter }, null, 1));
        if (empty) code = 4;   // fieldsAfter may still hold the answer
      }
    }
  }
} finally {
  await browser.close();
}
process.exit(code);
