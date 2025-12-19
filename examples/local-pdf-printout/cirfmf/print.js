import './filereader.js';
import { generateInvoice, generatePDFUPO } from './ksef-fe-invoice-converter.js';
import fs from 'fs';

// process.argv description:
// [0] => node binary
// [1] => print.js location
// [2] => mode (either invoice or UPO)
// [3] => source file
// [4] => destination file (pdf)
// [5] => KSeF invoice number (only for invoice mode)
// [6] => KSeF invoice QRCode (only for invoice mode)

const mode = process.argv[2].toLowerCase();

if (mode == "invoice") {
    if (process.argv.length < 7) {
        throw new Error("usage: node print.js invoice <source file> <destination file> <KSeFRefNo> <KSeFQRCodeURL>")
    }
    const additionalData = {
        nrKSeF: process.argv[5],
        qrCode: process.argv[6]
    };
    generateInvoice(process.argv[3], additionalData, 'blob').then((data) => {
        data.arrayBuffer().then(
            (arrayBuffer) => {
                const buffer = Buffer.from(arrayBuffer);
                fs.writeFile(process.argv[4], buffer, err => {
                    if (err) console.error(err);
                });
            }
        ).catch((err) => console.error(err));
    })
} else if (mode == "upo") {
    generatePDFUPO(process.argv[3]).then((data) => {
        data.arrayBuffer().then((arrayBuffer) => {
            const buffer = Buffer.from(arrayBuffer);
            fs.writeFile(process.argv[4], buffer, err => {
                if (err) console.error(err);
            })
        })
    });
} else {
    throw new Error("unsupported mode: " + mode)
}
