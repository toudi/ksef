<script setup>
import { QRCodeSVG } from '@akamfoad/qrcode';
import Invoice from "./components/Invoice.vue";
import { xml2object } from "@/xml2object";
import { annotateInvoice } from "@/annotate-invoice";


const qrCodeURL = document.querySelector("meta[name=\"invoice:qrcode\"]").content
const qrcode = new QRCodeSVG(qrCodeURL);
const invoiceXml = atob(document.querySelector("meta[name=invoice]").content);
const invoice = xml2object(invoiceXml);

if (!Array.isArray(invoice.Faktura.Fa.FaWiersz)) {
  invoice.Faktura.Fa.FaWiersz = [invoice.Faktura.Fa.FaWiersz]
}
const context = {
  invoice: invoice,
  annotations: annotateInvoice(invoice),
  qrcode: {
    url: qrCodeURL,
    svg: qrcode.toString(),
  },
  seiRefNo: document.querySelector("meta[name=\"invoice:seiRefNo\"]").content,
}
</script>

<template>
  <Invoice :context="context" />
</template>
