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
const output = fs.createWriteStream('./points.txt')

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

    // let timestamp = date.valueOf() * 1000000
    let timestamp = date.valueOf()
    const tagStr = tags.join(',')

    if (data.details && data.details.length > 0) {
      data.details.forEach(p => {
        timestamp = timestamp + 1

        const values = []

        values.push(`latency=${p.latency}`)

        const error = JSON.stringify(p.error)
        values.push(`error=${error}`)

        const status = JSON.stringify(p.status)
        values.push(`status=${status}`)

        const valStr = values.join(',')

        const line = `${measurementName},${tagStr} ${valStr} ${timestamp}000000`

        output.write(line + os.EOL)
      })
    }
  } catch (e) {
    console.log(e)
  }
})

output.end()
