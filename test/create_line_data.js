#!/usr/bin/env node

const fs = require('fs')
const path = require('path')
const os = require('os')

const files = fs.readdirSync(__dirname)

const runFiles = files.map(p => {
  const ext = path.extname(p)
  const bn = path.basename(p)
  if (ext === '.json' && bn.indexOf('run') >= 0) {
    const inputPath = path.join(__dirname, p)
    return inputPath
  }
})

let n = 0
const lines = []
runFiles.forEach(rf => {
  try {
    if (!rf) {
      return
    }

    n++
    if (n > 27) {
      console.log('maximum reached skipping...')
      return
    }

    const content = fs.readFileSync(rf, 'utf8')
    const data = JSON.parse(content)
    const date = new Date(data.date)

    const measurementName = 'ghz_run'

    const tags = []

    for (let k in data.options) {
      let v = data.options[k]

      if (typeof v === 'object') {
        v = JSON.stringify(data.options[k])
        // escape quotes
        v = JSON.stringify(v)
      } else {
        v = JSON.stringify(data.options[k])
      }

      // v = JSON.stringify(data.options[k])

      tags.push(`${k}=${v}`)
    }

    const values = []

    values.push(`count=${data.count}`)
    values.push(`total=${data.total}`)
    values.push(`average=${data.average}`)
    values.push(`fastest=${data.fastest}`)
    values.push(`slowest=${data.slowest}`)
    const rps = Number.parseFloat(data.rps).toFixed(2)
    values.push(`rps=${rps}`)

    if (data.latencyDistribution && data.latencyDistribution.length > 0) {
      const median = data.latencyDistribution.find(e => {
        return e.percentage === 50
      })

      if (median && median.latency) {
        values.push(`median=${median.latency}`)
      }

      const nine5 = data.latencyDistribution.find(e => {
        return e.percentage === 95
      })

      if (nine5 && nine5.latency) {
        values.push(`p95=${nine5.latency}`)
      }
    }

    let errorCount = 0

    for (let k in data.errorDistribution) {
      let v = data.errorDistribution[k]
      errorCount += v
    }

    values.push(`errors=${errorCount}`)

    const hasErrors = errorCount ? 'true' : 'false'
    tags.push(`hasErrors=${hasErrors}`)

    const dateISO = date.toISOString()
    tags.push(`date="${dateISO}"`)

    const timestamp = date.valueOf() * 1000000
    const tagStr = tags.join(',')
    const valStr = values.join(',')

    lines.push(`${measurementName},${tagStr} ${valStr} ${timestamp}`)
  } catch (e) {
    console.log(e)
  }
})

if (lines && lines.length > 0) {
  const output = fs.createWriteStream('./lines.txt')
  lines.forEach(l => {
    output.write(l + os.EOL)
  })
  output.end()
}
