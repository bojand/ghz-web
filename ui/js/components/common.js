function resultPass (test, value, type) {
  if (!test) {
    return false
  }

  if (type === 'fastest' || type === 'slowest' || type === 'mean') {
    if (
      test.thresholds &&
      test.thresholds[type] &&
      test.thresholds[type].threshold > value
    ) {
      return true
    }

    return false
  }

  if (type === 'RPS') {
    if (
      test.thresholds &&
      test.thresholds[type] &&
      test.thresholds[type].threshold < value
    ) {
      return true
    }

    return false
  }

  return false
}

function classifyResult (test, value, type) {
  if (!test) {
    return ''
  }

  const pass = resultPass(test, value, type)

  return pass ? 'is-success' : 'is-warning'
}

function iconifyResult (test, value, type) {
  if (!test) {
    return ''
  }

  const pass = resultPass(test, value, type)

  return pass ? 'checkbox-marked-circle-outline' : 'alert-circle-outline'
}

function formatFloat (val, fixed) {
  if (!Number.isInteger(fixed)) {
    fixed = 2
  }

  return Number.parseFloat(val).toFixed(fixed)
}

function formatMs (val) {
  return Number.parseFloat(val / 1000).toFixed(2)
}

function formatNano (val) {
  return Number.parseFloat(val / 1000000).toFixed(2)
}

export default {
  methods: {
    resultPass,
    classifyResult,
    iconifyResult,
    formatMs,
    formatNano,
    formatFloat
  }
}
