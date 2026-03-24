declare module 'virtual:pwa-info' {
	export interface PwaInfo {
		webManifest: {
			linkTag: string;
		};
	}
	export const pwaInfo: PwaInfo | undefined;
}
