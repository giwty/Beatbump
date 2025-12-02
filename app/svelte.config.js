import adapterCf from "@sveltejs/adapter-cloudflare";
import adapterCfw from "@sveltejs/adapter-cloudflare-workers";
import netlify from "@sveltejs/adapter-netlify";
import node from "@sveltejs/adapter-node";
import staticadapter from "@sveltejs/adapter-static";
import vercel from "@sveltejs/adapter-vercel";
import path from "path";
import sveltePreprocess from "svelte-preprocess";
import { fileURLToPath } from "url";


const dev = process.env["NODE_ENV"] === "development";
const ENV_ADAPTER = process.env["BB_ADAPTER"] ?? "staticadapter";
const ENV_IFRAME = process.env["ALLOW_IFRAME"] ?? "false";

const adapters = {
    "cloudflare-workers": adapterCfw({}),
    cloudflare: adapterCf(),
    vercel: vercel({}),
    netlify: netlify({}),
    node: node({ "precompress": false }),
    staticadapter: staticadapter({
        pages: 'build',
        assets: 'build',
        fallback: 'index.html',
        precompress: false,
        strict: true
    })
};

const adapter = adapters[ENV_ADAPTER];
/** @type {import('@sveltejs/kit').Config} */
const config = {
    preprocess: sveltePreprocess({
        sass: false,

        scss: {
            includePaths: ["./src/"],

            prependData: `@use "global/redesign/base/_variables.scss" as *;`, stripIndent: true,
        },
        postcss: { configFilePath: path.resolve("./postcss.config.cjs") },
    }),

    kit: {
        adapter: dev ? node() : adapters["staticadapter"], env: { dir: '../' },
        alias: {
            $stores: path.resolve("./src/lib/stores"),
            $api: path.resolve("./src/routes/(app)/api/_lib"),
            $components: path.resolve("./src/lib/components"),
        },
        prerender: { concurrency: 3 },
        files: {
            assets: "static",
            lib: "src/lib",
            routes: "src/routes",
            //  serviceWorker: "src/service-worker",
            appTemplate: "src/app.html",
            // hooks: {server: "src/hooks.server"},
        },
        csp: {
            directives: {
                'frame-ancestors': ['self', ENV_IFRAME ? '*' : '']
            }
        },
        version: { pollInterval: 600000 },
    },
    onwarn(warning, defaultHandler) {
        if (warning.code === "css-unused-selector") return;

        defaultHandler(warning);
    },
};
export default config;
