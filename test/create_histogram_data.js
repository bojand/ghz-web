#!/usr/bin/env node

const fs = require('fs')
const path = require('path')
const os = require('os')

const files = fs.readdirSync(__dirname)

const rf = files.find(p => {
  const ext = path.extname(p)
  const bn = path.basename(p)
  return ext === '.json' && bn.indexOf('run') >= 0
})

let n = 0
const lines = []
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

  const measurementName = 'ghz_histogram'

  const tags = []

  for (let k in data.options) {
    let v = data.options[k]

    if (typeof v === 'object') {
      // v = JSON.stringify(data.options[k])
      // // escape quotes
      // v = JSON.stringify(v)
      continue
    } else {
      v = JSON.stringify(data.options[k])
    }

    // v = JSON.stringify(data.options[k])

    tags.push(`${k}=${v}`)
  }

  let timestamp = date.valueOf()

  if (data.histogram && data.histogram.length > 0) {
    data.histogram.forEach(b => {
      for (let c = 0; c < b.count; c++) {
        const values = []
        timestamp += 1
        // values.push(`${b.mark}=${b.count}`)

        values.push(`mark=${b.mark}`)

        const tagStr = tags.join(',')
        const valStr = values.join(',')

        lines.push(`${measurementName},${tagStr} ${valStr} ${timestamp}000000`)
      }
    })
  }
} catch (e) {
  console.log(e)
}

if (lines && lines.length > 0) {
  const output = fs.createWriteStream('./histogram.txt')
  lines.forEach(l => {
    output.write(l + os.EOL)
  })
  output.end()
}
