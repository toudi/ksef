export const vatRateASString = (vatRate) => {
    if (parseFloat(vatRate) !== NaN) {
        return vatRate + " %"
    }
    return vatRate;
}
