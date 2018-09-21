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

    const measurementName = 'ghz_histogram'

    const tags = []

    // for (let k in data.options) {
    //   let v = data.options[k]

    //   if (typeof v === 'object') {
    //     // v = JSON.stringify(data.options[k])
    //     // // escape quotes
    //     // v = JSON.stringify(v)
    //     continue
    //   } else {
    //     v = JSON.stringify(data.options[k])
    //   }

    //   // v = JSON.stringify(data.options[k])

    //   tags.push(`${k}=${v}`)
    // }

    let timestamp = date.valueOf()

    const dateISO = date.toISOString()

    if (data.histogram && data.histogram.length > 0) {
      data.histogram.forEach(b => {
        timestamp += 1

        const values = []
        values.push(`count=${b.count}`)

        const curTags = tags.slice()
        curTags.push(`date="${dateISO}"`)
        curTags.push(`mark_s=${b.mark}`)
        curTags.push(`mark_ms=${b.mark * 1000}`)
        curTags.push(`mark_ns=${b.mark * 1000000000}`)
        
        const tagStr = curTags.join(',')
        const valStr = values.join(',')

        lines.push(`${measurementName},${tagStr} ${valStr} ${timestamp}000000`)
      })
    }

    // tags.push(`date="${dateISO}"`)

    // const values = []

    // if (data.histogram && data.histogram.length > 0) {
    //   let bn = 0

    //   data.histogram.forEach(b => {
    //     tags.push(`bucket_${bn}_mark_s=${b.mark}`)
    //     tags.push(`bucket_${bn}_mark_ms=${b.mark * 1000}`)
    //     tags.push(`bucket_${bn}_mark_ns=${b.mark * 1000000000}`)

    //     values.push(`bucket_${bn}_count=${b.count}`)
    //     bn = bn + 1
    //   })
    // }

    // const tagStr = tags.join(',')
    // const valStr = values.join(',')

    // lines.push(`${measurementName},${tagStr} ${valStr} ${timestamp}000000`)
  } catch (e) {
    console.log(e)
  }
})

if (lines && lines.length > 0) {
  const output = fs.createWriteStream('./histogram.txt')
  lines.forEach(l => {
    output.write(l + os.EOL)
  })
  output.end()
}
