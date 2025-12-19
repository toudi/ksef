# Important info

In order for this script to work, you first have to build the ministry's library and hope it does not change in the future :D

Here's a quick tutorial.

1. Clone the repository:

https://github.com/CIRFMF/ksef-pdf-generator

2. Go into the directory

```
cd ksef-pdf-generator
```

3. install requirements and build the library:

```
npm install
npm run build
```

4. Copy `ksef-fe-invoice-converter.js` here (i.e. to the same location where you're using the helper scripts)

5. That's it. You can now use the helper script like so:

## invoice mode

```
node print.js invoice invoice-source.xml invoice-dest.pdf 5555555555-20250808-9231003CA67B-BE https://ksef-test.mf.gov.pl/client-app/invoice/5265877635/26-10-2025/HS5E1zrA8WVjDNq_xMVIN5SD6nyRymmQ-BcYHReUAa0

```

## upo mode

```
node print.js upo upo-source.xml upo-dest.pdf
```
