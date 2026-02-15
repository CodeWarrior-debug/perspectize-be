import type { Metric } from 'web-vitals';

function sendToAnalytics(metric: Metric) {
	console.debug('[Web Vitals]', metric.name, Math.round(metric.value), 'ms', {
		id: metric.id,
		rating: metric.rating
	});
}

export function reportWebVitals() {
	import('web-vitals').then(({ onCLS, onINP, onLCP, onFCP, onTTFB }) => {
		onCLS(sendToAnalytics);
		onINP(sendToAnalytics);
		onLCP(sendToAnalytics);
		onFCP(sendToAnalytics);
		onTTFB(sendToAnalytics);
	});
}
