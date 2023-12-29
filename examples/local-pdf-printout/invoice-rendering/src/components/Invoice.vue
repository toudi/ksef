<template>
    <div class="container-fluid">
        <div class="row">
            <div class="col">
                <div id="seiQRCode" v-html="context.qrcode.svg"></div>
            </div>
            <div class="col text-end align-self-center">
                <div class="text-end">
                    <h1>Faktura VAT</h1>
                    <div class="d-flex justify-content-end">
                        <table>
                            <tr>
                                <td>Numer faktury</td>
                                <td>{{ invoice.Fa.P_2 }}</td>
                            </tr>
                            <tr>
                                <td>Data wystawienia</td>
                                <td>{{ invoice.Fa.P_1 }}</td>
                            </tr>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        <div class="row my-4 text-center">
            <p>Numer faktury w KSeF: <a target="_blank" :href="context.qrcode.url"><span class="fw-bold">{{ context.seiRefNo
            }}</span></a>
            </p>
        </div>
        <div class="row my-4">
            <div class="col">
                <subject :subject="invoice.Podmiot1" />
            </div>
            <div class="col">
                <subject :subject="invoice.Podmiot2" />
            </div>

            <div class="col" v-if="invoice.Podmiot3 !== undefined">
                <subject :subject="invoice.Podmiot3" />
            </div>
        </div>
        <div class="row my-4">
            <div class="col">
                <items :items="annotations.items" :total="annotations.total" />
            </div>
        </div>
        <div class="row my-4">
            <div class="col">
                <p>Podsumowanie według stawek VAT</p>
                <table class="table">
                    <thead>
                        <tr>
                            <th>Stawka</th>
                            <th>Netto</th>
                            <th>VAT</th>
                            <th>Brutto</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="total, rate of annotations.total.by_rate" :key="rate">
                            <td>{{ vatRateASString(rate) }}</td>
                            <td>{{ total.net.toFixed(2) }}</td>
                            <td>{{ total.vat.toFixed(2) }}</td>
                            <td>{{ total.gross.toFixed(2) }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import subject from "./subject.vue"
import items from "./items.vue";
import { vatRateASString } from '@/utils';

const props = defineProps(["context"]);
const context = ref(props.context);
const invoice = ref(props.context.invoice.Faktura);
const annotations = ref(props.context.annotations);
</script>

<style>
svg {
    height: 100%;
}

#seiQRCode {
    height: 15em;
}
</style>