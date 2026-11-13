// calculator_net_map.mjs: map the calculator links on calculator.net's index pages.
//
// Usage: node scripts/calculator_net_map.mjs <robots.txt URL> <index URL>...
// Loads robots.txt first, then each index page in order, one page at a time,
// waiting a random 4-8 s between loads. Prints:
//   ROBOTS<TAB><status>       then robots.txt verbatim between ---ROBOTS--- markers
//   STATUS<TAB><url><TAB><http status>
//   ROW<TAB><index url><TAB><section><TAB><name><TAB><absolute url>
// Only anchor text and href are collected. A 403/429 or block page stops the run.

import { chromium } from '@playwright/test';

const BASE = 'https://www.calculator.net/';
const [robotsUrl, ...indexUrls] = process.argv.slice(2);
if (!robotsUrl || indexUrls.length === 0) {
  console.error('usage: calculator_net_map.mjs <robots.txt URL> <index URL>...');
  process.exit(2);
}

const pause = () => new Promise(r => setTimeout(r, 4000 + Math.random() * 4000));

// disallowed reports whether robots.txt (User-agent: * group) disallows path.
function disallowed(robots, path) {
  let applies = false;
  const rules = [];
  for (const raw of robots.split(/\r?\n/)) {
    const line = raw.replace(/#.*/, '').trim();
    const m = line.match(/^([A-Za-z-]+)\s*:\s*(.*)$/);
    if (!m) continue;
    const key = m[1].toLowerCase();
    if (key === 'user-agent') applies = m[2].trim() === '*';
    else if (applies && key === 'disallow' && m[2]) rules.push(m[2].trim());
  }
  return rules.some(r => path.startsWith(r));
}

const browser = await chromium.launch({ headless: true });
let exitCode = 0;
try {
  const page = await browser.newPage();
  const resp = await page.goto(robotsUrl, { waitUntil: 'domcontentloaded', timeout: 30000 });
  const robots = resp ? await resp.text() : '';
  console.log(`ROBOTS\t${resp ? resp.status() : 'none'}`);
  console.log('---ROBOTS---');
  process.stdout.write(robots);
  console.log('\n---END ROBOTS---');

  for (const url of indexUrls) {
    const path = new URL(url).pathname;
    if (disallowed(robots, path)) {
      console.log(`STATUS\t${url}\tSKIPPED (robots.txt disallows ${path})`);
      continue;
    }
    await pause();
    const r = await page.goto(url, { waitUntil: 'domcontentloaded', timeout: 30000 });
    const status = r ? r.status() : 0;
    console.log(`STATUS\t${url}\t${status}`);
    const title = (await page.title()).toLowerCase();
    if (status === 403 || status === 429 || /access denied|attention required|just a moment|captcha/.test(title)) {
      console.log(`BLOCKED\t${url}\t${status}\t${title}`);
      exitCode = 3;
      break;
    }
    // Links inside the main content column, each with the nearest heading above it.
    const rows = await page.evaluate(() => {
      const root = document.querySelector('#contentout') || document.querySelector('#content') || document.body;
      const out = [];
      let section = '';
      const walker = document.createTreeWalker(root, NodeFilter.SHOW_ELEMENT);
      for (let el = walker.currentNode; el; el = walker.nextNode()) {
        if (/^H[1-4]$/.test(el.tagName)) section = el.textContent.trim();
        else if (el.tagName === 'A' && el.getAttribute('href')) {
          out.push({ section, name: el.textContent.trim(), href: el.getAttribute('href') });
        }
      }
      return out;
    });
    for (const row of rows) {
      const abs = new URL(row.href, BASE).href;
      console.log(`ROW\t${url}\t${row.section}\t${row.name}\t${abs}`);
    }
  }
} finally {
  await browser.close();
}
process.exit(exitCode);
