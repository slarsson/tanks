'use strict';

const prod = (process.env.NODE_ENV === 'production'); 

const path = require('path');
const gulp = require('gulp');
const webpack = require('webpack');
const Browser = require('browser-sync');
const htmlmin = require('gulp-htmlmin');
const del = require('del');
const sass = require('gulp-sass'); 
	  sass.compiler = require('node-sass');

let browser;
if (!prod) {
	browser = Browser.create();
}

const webpackConfig = {
	mode: prod ? 'production' : 'development',
	entry: './src/index.js',
	module: {
		rules: [
			{
		  		test: /\.tsx?$/,
		  		use: 'ts-loader',
		  		exclude: /node_modules/
			}
	  	]
	},
	resolve: {
		extensions: ['.tsx', '.ts', '.js']
	},
	output: {
		filename: 'bundle.js',
		path: path.resolve(__dirname, 'build')
	},
	plugins: [
		new webpack.DefinePlugin({
			WS_HOST: process.env.WS_HOST ? '"'+process.env.WS_HOST+'"' : '"localhost:1337"'
		})
	]
};

const wp = () => {
    return new Promise((resolve, reject) => webpack(webpackConfig, (err, stats) => {
        if (err) {
			reject(err);
			return;
		}
		console.log(stats.toString({colors: true}));
		resolve()
    }));
};

const html = () => {
	return gulp.src('*.html')
    	.pipe(htmlmin({collapseWhitespace: true}))
        .pipe(gulp.dest('build'));
};

const css = () => {
    return gulp.src('./*.scss')
        .pipe(sass({outputStyle: 'compressed'}).on('error', sass.logError))
        .pipe(gulp.dest('./build'));
};

const reload = (cb) => {
    browser.reload();
    cb();
};

const clean = (cb) => {
    del.sync(['build']);
	cb();
};

const copy = () => {
	return gulp.src('src/assets/**/*')
		.pipe(gulp.dest('build/assets'));
};

const wasm = () => {
	return gulp.src('assets/*')
		.pipe(gulp.dest('build'));
};

const server = () => {
	browser.init({
        server: './build',
        open: false,
	});
	
	gulp.watch('./assets/*', gulp.series('wasm', 'reload'));

	//gulp.watch('./src/assets/*', gulp.series('delete', 'copy', 'reload'));
    gulp.watch('./src/*', gulp.series('webpack', 'reload'));
    //gulp.watch('./*.html', gulp.series('html', 'reload'));
    gulp.watch('./*.scss', gulp.series('css', 'reload'));
};

gulp.task('html', html);
gulp.task('webpack', wp);
gulp.task('css', css);
gulp.task('reload', reload);
gulp.task('delete', clean);
gulp.task('copy', copy);
gulp.task('server', server);
gulp.task('wasm', wasm);
gulp.task('build', gulp.series('delete', 'copy', 'html', 'css', 'webpack', 'wasm'));

if (prod) {
    gulp.task('default', gulp.series('build'));
} else {
    gulp.task('default', gulp.series('build', 'server'));
}