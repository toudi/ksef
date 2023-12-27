/**
 * Copyright 2017 Google Inc. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

'use strict';

const puppeteer = require('puppeteer-core');

/**
 * process.argv[2] => chrome location
 * process.argv[3] => source HTML file
 * process.argv[4] => destination PDF
 */

if(process.argv.length < 5) {
  process.exit(1);
}

(async () => {
  const browser = await puppeteer.launch({
    executablePath: process.argv[2],
    args: ["--disable-web-security"]
  });

  const page = await browser.newPage();

  await page.goto('file://' + process.argv[3], {
    waitUntil: 'networkidle2',
  });
  await page.emulateMediaType('print');

  // await page.emulateMedia('print');
  // page.pdf() is currently supported only in headless mode.
  // @see https://bugs.chromium.org/p/chromium/issues/detail?id=753118
  await page.pdf({
    path: process.argv[4],
    format: 'A4',
    preferCSSPageSize: true,
    printBackground: true,
    scale: 0.8,
    margin: {
      left: '1cm',
      right: '1cm',
      bottom: '1cm',
      top: '1cm',
    }
  });

  await browser.close();
})();
