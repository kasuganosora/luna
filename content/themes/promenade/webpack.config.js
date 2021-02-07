const path = require("path");
const MiniCssExtractPlugin = require('mini-css-extract-plugin');


module.exports = {
    entry: {
        app: "./assets/sass/styles.js"
    },
    output: {
        publicPath: __dirname + "/assets/",
        path: path.resolve(__dirname, "assets"),
        filename: "[name].bundle.js"
    },
    plugins: [
        new MiniCssExtractPlugin({
            filename: "css/style.css",
        }),
    ],
    module: {
        rules: [
            {
                test: /\.scss$/,
                use: [
                    {
                        loader: MiniCssExtractPlugin.loader,
                    },
                    {
                        loader: "css-loader",
                        options: {url: false}
                    },
                    {
                        loader: "sass-loader",
                    }
                ]
            }]
    }
};