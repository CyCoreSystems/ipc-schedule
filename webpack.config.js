var webpack = require('webpack');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var LessPluginCleanCSS = require('less-plugin-clean-css');

module.exports = {
   entry: './app/index',
   output: {
      path: __dirname + '/public',
      publicPath: "/app/",
      filename: 'bundle.js'
   },
   plugins: [
      new webpack.ProvidePlugin({
         riot: 'riot'
      }),
      new ExtractTextPlugin('styles.css')
   ],
   devtool: 'source-map',
   lessLoader: {
      lessPlugins: [
         new LessPluginCleanCSS({advanced: true})
      ]
   },
   module: {
      preLoaders: [
         { test: /\.tag$/, exclude: /node_modules/, loader: 'riotjs-loader', query: { type: 'none' } }
      ],
      loaders: [
         {
            test: /\.js$|\.tag$/,
            exclude: /node_modules/,
            loader: 'babel',
            query: {
               presets: ['es2015']
            }
         },
         /*
         {
            test: /\.less$/,
            loader: ExtractTextPlugin.extract('css?sourceMap!less?sourceMap')
         },
         */
         {
            test: /\.less$/,
            loader: 'style!css?sourceMap!less?sourceMap!'
         }
      ]
   },
   devServer: {
      hot: true,
      inline: true,
      debug: true,
      progress: true,
      port: 3000,
      devtool: 'eval-source-map',
      open: true,
      watchOptions: {
         aggregateTimeout: 300,
         poll: 1000
      },
      stats: { colors: true },
      proxy: [{
         path: /^(?!.*\.hot-update\.js)(.*)$/,
         target: 'http://localhost:9000'
      }]
   }
};
