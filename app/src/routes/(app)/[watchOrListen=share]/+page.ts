import { redirect } from "@sveltejs/kit";
import {APIClient} from "$lib/api";
export const load = async ({ url, fetch }) => {
	const id =
		url.searchParams.get("id") ??
		url.searchParams.get("v") ??
		url.searchParams.get("videoId");
	const playlist = url.searchParams.get("list") || undefined;

	if (!id) {
		throw redirect(301, "/trending");
	}

	const [data, list] = await Promise.all([
        APIClient.fetch(
			`/api/v1/player.json?videoId=${id ? id : ""}${
				playlist ? `&playlistId=${playlist}` : ""
			}`,
		).then((res) => res.json()),
        APIClient.fetch(
			`/api/v1/next.json?videoId=${id ? id : ""}${
				playlist ? `&playlistId=${playlist}` : ""
			}`,
		).then((res) => res.json()),
	]);
	const {
		videoDetails: {
			title = "",
			videoId = "",
			thumbnail: { thumbnails = [] } = {},
		} = {},
	} = data;

    console.log(data)

	return {
		title,
		thumbnails,
		videoId,
		playlist,
		related: list,
		data,
	};
};
