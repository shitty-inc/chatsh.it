const webpack = require('webpack')

module.exports = function override(config, env) {
  config.resolve = {
    fallback: {
      buffer: false,
      crypto: false,
      fs: false,
      os: false,
      path: false,
      stream: false,
      util: false
    },
    extensions: [ '.go', '.tsx', '.ts', '.js' ]
  }

  config.module.rules = config.module.rules.map(rule => {
    if (rule.oneOf instanceof Array) {
      return {
        ...rule,
        oneOf: [
          {
            test: /\.go/,
            use: ['@fiedka/golang-wasm-async-loader']
          },
          ...rule.oneOf
        ]
      };
    }

    return rule;
  });

  config.plugins = (config.plugins || []).concat([
    new webpack.ProvidePlugin({
      Buffer: ['buffer', 'Buffer'],
      process: 'process/browser'
    })
  ])
  return config
}
