import axios from 'axios'

export default {
  project: null,
  test: null,
  run: null,

  async fetchProject (id) {
    const { data } = await axios.get(`http://localhost:3000/api/projects/${id}`)

    this.project = data
    return data
  },

  async updateProject (projectData) {
    const { data } = await axios.put(
      `http://localhost:3000/api/projects/${projectData.id}`,
      projectData
    )

    this.project = data
    return data
  },

  async fetchTest (projectId, testId) {
    const { data } = await axios.get(
      `http://localhost:3000/api/projects/${projectId}/tests/${testId}`
    )

    this.test = data
    return data
  },

  async updateTest (projectId, testData) {
    console.log(JSON.stringify(testData))
    const { data } = await axios.put(
      `http://localhost:3000/api/projects/${projectId}/tests/${testData.id}`,
      testData
    )

    this.test = data
    return data
  },

  async fetchLatestRun (projectId, testId) {
    const { data } = await axios.get(
      `http://localhost:3000/api/projects/${projectId}/tests/${testId}/runs/latest`
    )

    return data
  }
}
