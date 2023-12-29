import type { DefaultTheme } from "vitepress";
import vitepressConfig from "./config.mjs";
import { defineUserConfig } from "vitepress-export-pdf";

function extractLinksFromConfig(config: DefaultTheme.Config) {
	const links: string[] = [];

	function extractLinks(sidebar: DefaultTheme.SidebarItem) {
		if (Array.isArray(sidebar)) {
			for (const item of sidebar) {
				extractLinks(item)
			}
		}

		if (sidebar.items)
			extractLinks(sidebar.items);

		else if (sidebar.link)
			links.push(`${sidebar.link}.html`);
	}

	for (const key in config.sidebar) {
		extractLinks(config.sidebar[key]);
	}

	return links;
}

const routeOrder = extractLinksFromConfig(vitepressConfig.themeConfig!);

export default defineUserConfig({
	routePatterns: ["!/index.html"],
	pdfOptions: {
		printBackground: true,
		margin: {
			bottom: 60,
			left: 25,
			right: 25,
			top: 60,
		},
	},
	sorter: (pageA, pageB) => {
		const aIndex = routeOrder.findIndex(route => route === pageA.path);
		const bIndex = routeOrder.findIndex(route => route === pageB.path);
		return aIndex - bIndex;
	},
});