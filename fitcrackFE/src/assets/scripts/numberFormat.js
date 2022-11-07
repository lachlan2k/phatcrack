export default function (number, options = {}, locale = 'en-US') {
  return new Intl.NumberFormat(locale, options).format(parseFloat(number))
}