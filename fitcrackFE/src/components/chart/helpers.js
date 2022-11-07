import pattern from 'patternomaly'
import moment from 'moment'

const gridColor = '#9997'
export { gridColor }

function* colorGenerator () {
  let hue = 24 * 8
  while (true) {
    hue += 24 * 11 // 24deg base step = 15 colors
    yield `hsl(${hue}, 70%, 55%)`
  }
}
let gen = colorGenerator()
let colors = []
let patterns = []

getColors(15)
patterns = [...colors]

export function getColors (count, usePatterns = false) {
  const missingColors = count - colors.length
  if (missingColors > 0) {
    let newColors = []
    for (let i = 0; i < missingColors; i++) {
      newColors.push(gen.next().value)
    }
    const newPatterns = pattern.generate(newColors)
    colors = [...colors, ...newColors]
    patterns = [...patterns, ...newPatterns]
  }
  if (usePatterns) {
    return patterns
  } else {
    return colors
  }
}

export function prepareLines (data) {
  const colors = getColors(data.length)
  return data.map((set, i) => {
    set.backgroundColor = `hsla(${colors[i].substring(4, colors[i].length - 1)}, .4)`
    set.pointBackgroundColor = colors[i]
    set.borderColor = colors[i]
    return set
  })
}

export function localizeTime (data) {
  return data.map(set => {
    let localized = {}
    for (const [time, value] of Object.entries(set.data)) {
      localized[moment.utc(time).local()] = value
    }
    set.data = localized
    return set
  })
}

export function plotUsage (usageData, metric) {
  return take(usageData, 180)
    .reduce((data, point) => ({
      ...data,
      [moment.utc(point.time).local()]: point[metric]
    }), {})
}

function take (array, amount) {
  if (array.length <= amount) return array
  const step = Math.floor(array.length / amount)
  let result = []
  for (let i = array.length - 1; i > 0; i -= step) {
    if (i < 0) break
    result.push(array[i])
  }
  return result
}