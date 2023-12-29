<template>
    <div class="card">
        <div class="card-header">Sprzedawca</div>
        <div class="card-body">
            <strong>{{ subject.DaneIdentyfikacyjne.Nazwa }}</strong>
            <p>{{ subject.Adres.AdresL1 }}</p>
            <p>{{ subject.Adres.AdresL2 }}</p>
            <p>NIP: {{ subject.DaneIdentyfikacyjne.NIP }}</p>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
const props = defineProps(["subject"]);
const subject = ref(props.subject);
</script>