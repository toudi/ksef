const InvoiceItemRenameFields = {
    "P_7": "name",
    "P_8A": "units",
    "P_8B": "quantity",
    "P_9A": "unit-price-net",
    "P_9B": "unit-price-gross",
    "P_12": "vat-rate",
}

const round = (value) => {
    return Math.round(value * 100) / 100
}

const parseTotals = (item) => {
    let total = {
        net: 0.0,
        gross: 0.0,
        vat: 0.0,
    }

    const vatQuantizer = 1 + item["vat-value"]

    if (item["unit-price-net"] !== undefined) {
        // calculate from net to gross
        total.net = round(item["unit-price-net"] * item["quantity"])
        total.gross = round(total.net * vatQuantizer)
    } else {
        total.gross = round(item["unit-price-gross"] * item["quantity"])
        total.net = round(total.gross / vatQuantizer)
    }
    
    total.vat = round(total.gross - total.net);

    return total;
}

export const annotateInvoice = (invoice) => {
    let total = {
        by_rate: {},
        total: {},
    }
    let items = []
    let vat_rate = ""
    let totalTemplate = {
        net: 0.0,
        vat: 0.0,
        gross: 0.0,
    }

    total["total"] = {...totalTemplate}

    for (let item of invoice.Faktura.Fa.FaWiersz) {
        for (const [old_key, new_key] of Object.entries(InvoiceItemRenameFields)) {
            if (item[old_key] !== undefined) {
                item[new_key] = item[old_key]
                delete(item[old_key])
            }
        }
        item["vat-value"] = 0.0
        if (typeof item["vat-rate"] === "number") {
            item["vat-value"] = item["vat-rate"] / 100;
        }
        items.push(item);
        vat_rate = item["vat-rate"]
        if (total["by_rate"][vat_rate] === undefined) {
            total["by_rate"][vat_rate] = {...totalTemplate}
        }
        var totals = parseTotals(item);
        item["amount"] = totals;

        total["by_rate"][vat_rate].net = round(totals.net + total["by_rate"][vat_rate].net);
        total["by_rate"][vat_rate].gross = round(totals.gross + total["by_rate"][vat_rate].gross);
        total["by_rate"][vat_rate].vat = round(totals.vat + total["by_rate"][vat_rate].vat)

        total["total"].net = round(totals.net + total["total"].net);
        total["total"].gross = round(totals.gross + total["total"].gross);
        total["total"].vat = round(totals.vat + total["total"].vat);

    }
    return {
        items: items,
        total: total,
    }
}