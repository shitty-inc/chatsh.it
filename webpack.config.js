const webpack = require('webpack');
const path = require('path');

module.exports = {
  entry: './src/app.ts',
  devtool: 'inline-source-map',
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader',
        exclude: /node_modules/
      },
      {
        test: /\.go/,
        use: ['@fiedka/golang-wasm-async-loader']
      }
    ]
  },
  resolve: {
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
  },
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  },
  plugins: [
    new webpack.ProvidePlugin({
      process: 'process/browser.js',
    }),
    new webpack.ProvidePlugin({
      Buffer: [require.resolve("buffer/"), "Buffer"],
    }),
  ]
};
