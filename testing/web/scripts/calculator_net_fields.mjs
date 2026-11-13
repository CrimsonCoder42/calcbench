// calculator_net_fields.mjs: list the form fields of one calculator page.
//
// Usage: node scripts/calculator_net_fields.mjs <calculator URL>
// One page load. Prints JSON: {url, status, fields: [{form, name, type, label, value, unit, options}]}.
// label is the text of the table cell or <label> before the field; unit is the short text right after it;
// options lists a select's or radio group's values. Nothing else from the page is read.
// Exit 3 on 403/429 or a block page.

import { chromium } from '@playwright/test';

const url = process.argv[2];
if (!url) {
  console.error('usage: calculator_net_fields.mjs <calculator URL>');
  process.exit(2);
}

const browser = await chromium.launch({ headless: true });
let code = 0;
try {
  const page = await browser.newPage();
  const resp = await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });
  // Some calculators draw their form from script after the first paint.
  await page.waitForLoadState("load", { timeout: 15000 }).catch(() => {});
  await page.waitForTimeout(1500);
  const status = resp ? resp.status() : 0;
  const title = (await page.title()).toLowerCase();
  if (status === 403 || status === 429 || /access denied|attention required|just a moment|captcha/.test(title)) {
    console.log(JSON.stringify({ url, status, blocked: true, title }));
    code = 3;
  } else {
    const fields = await page.evaluate(() => {
      const clean = s => (s || '').replace(/\s+/g, ' ').trim();
      const out = [];
      const seen = new Set();
      // Every field on the page, in or out of a form (some calculators keep their inputs outside any form).
      const els = document.querySelectorAll("input, select, textarea");
      for (const el of els) {
        const type = el.tagName === 'SELECT' ? 'select' : (el.type || el.tagName.toLowerCase());
        if (['hidden', 'submit', 'button', 'image', 'reset'].includes(type)) continue;
        const name = el.name || el.id;
        if (!name) continue;
        const form = el.form ? (el.form.name || el.form.id || '') : '';
        const key = form + '|' + name;
        if (type === 'radio') {
          const label = clean(el.labels && el.labels[0] ? el.labels[0].textContent : el.parentElement.textContent);
          const prev = out.find(f => f.form === form && f.name === name);
          const opt = { value: el.value, label, checked: el.checked };
          if (prev) { prev.options.push(opt); if (el.checked) prev.value = el.value; continue; }
          const td = el.closest('td');
          const row = el.closest('tr');
          const rowLabel = row && row.cells && row.cells[0] && row.cells[0] !== td ? clean(row.cells[0].textContent) : '';
          out.push({ form, name, type, label: rowLabel, value: el.checked ? el.value : '', unit: '', options: [opt] });
          seen.add(key);
          continue;
        }
        if (seen.has(key)) continue;
        seen.add(key);
        let label = '';
        if (el.labels && el.labels[0]) label = clean(el.labels[0].textContent);
        if (!label) {
          const td = el.closest('td');
          let prev = td && td.previousElementSibling;
          while (prev && !clean(prev.textContent)) prev = prev.previousElementSibling;
          if (prev) label = clean(prev.textContent);
        }
        let unit = '';
        const next = el.nextSibling;
        if (next && next.nodeType === 3) unit = clean(next.textContent);
        if (!unit) {
          const td = el.closest('td');
          const nx = td && td.nextElementSibling;
          if (nx && !nx.querySelector('input,select')) unit = clean(nx.textContent).slice(0, 30);
        }
        const f = { form, name, id: el.id || "", type, label: label.slice(0, 80), value: type === 'checkbox' ? el.checked : el.value, unit };
        if (type === 'select') f.options = [...el.options].map(o => ({ value: o.value, label: clean(o.textContent), selected: o.selected }));
        out.push(f);
      }
      return out;
    });
    console.log(JSON.stringify({ url, status, fields }, null, 1));
  }
} finally {
  await browser.close();
}
process.exit(code);
