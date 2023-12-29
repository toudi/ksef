<template>
    <table class="table">
        <thead>
            <tr>
                <th>L.P.</th>
                <th>Nazwa</th>
                <th>Jednostka</th>
                <th>Ilość</th>
                <th>Cena jednostkowa netto</th>
                <th>Cena jednostkowa bruto</th>
                <th>Stawka VAT</th>
                <th>Wartość sprzedaży netto</th>
                <th>Wartość VAT</th>
                <th>Wartość sprzedaży brutto</th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="item, i in items" :key="i">
                <td>
                    {{ item.NrWierszaFa }}
                </td>
                <td>{{ item.name }}</td>
                <td>{{ item.units }}</td>
                <td>{{ item.quantity }}</td>
                <td>
                    <template v-if="item['unit-price-net'] !== undefined">
                        {{ item["unit-price-net"] }}
                    </template>
                </td>
                <td>
                    <template v-if="item['unit-price-gross'] !== undefined">
                        {{ item["unit-price-gross"] }}
                    </template>
                </td>
                <td>{{ vatRateASString(item["vat-rate"]) }}</td>
                <td>{{ item["amount"].net.toFixed(2) }}</td>
                <td>{{ item["amount"].vat.toFixed(2) }}</td>
                <td>{{ item["amount"].gross.toFixed(2) }}</td>
            </tr>
        </tbody>
        <tfoot>
            <tr>
                <td colspan="7">RAZEM</td>
                <td>{{ total.total.net.toFixed(2) }}</td>
                <td>{{ total.total.vat.toFixed(2) }}</td>
                <td>{{ total.total.gross.toFixed(2) }}</td>
            </tr>
        </tfoot>
    </table>
</template>

<script setup>
import { ref } from 'vue';
import { vatRateASString } from '@/utils';
const props = defineProps(['items', 'total']);
const items = ref(props.items);
</script>