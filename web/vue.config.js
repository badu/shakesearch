module.exports = {
    filenameHashing: false,
    chainWebpack: config => {
        config
            .plugin('html')
            .tap(args => {
                args[0].title = "Shake Search"
                return args
            })
    },
}