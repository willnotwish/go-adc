const path = require('path')

module.exports = {
  chainWebpack: config => {
    const apiClient = process.env.VUE_APP_API_CLIENT // mock or server
    config.resolve.alias.set(
      'api-client',
      path.resolve(__dirname, `src/api/${apiClient}`)
    )
  },
  devServer: {
    proxy: {
      '/api': {
        target: 'http://192.168.0.48'
      },
      '/ws': {
        target: 'ws://192.168.0.48',
        ws: true
      }
    }
  }
}