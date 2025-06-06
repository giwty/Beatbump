import { goto } from "$app/navigation";
import type { IListItemRenderer } from "$lib/types/musicListItemRenderer";
import list from "$stores/list";

export const browseHandler = (pageType: string, browseId: string): void => {
	if (pageType.includes("ARTIST")) {
		goto(`/artist/${browseId}`);
	} else {
		pageType.includes("PLAYLIST")
			? goto("/playlist/" + browseId)
			: goto(
					"/release?type=" +
						encodeURIComponent(pageType) +
						"&id=" +
						encodeURIComponent(browseId),
			  );
	}
};
export async function clickHandler({
	item,
	index,
	isBrowseEndpoint,
	type,
	kind,
}: {
	item: IListItemRenderer;
	index: number;
	isBrowseEndpoint: boolean;
	type: string;
	kind: string;
}) {
    console.log(item)
	if (item.endpoint && item.endpoint?.pageType && item.endpoint?.browseId) {
		browseHandler(item.endpoint.pageType, item.endpoint?.browseId);
		return;
	}
	// console.log(item);
	if (type === "trending") {
		if (item.endpoint?.pageType?.match(/ALBUM|SINGLE/m)) {
			goto(
				"/release?type=" +
					encodeURIComponent(item?.endpoint?.pageType) +
					"&id=" +
					encodeURIComponent(item.endpoint?.browseId),
			);
		} else {
            console.log("initAutoMixSession")
			await list.initAutoMixSession({
                clickedItem: item,
				videoId: item.videoId,
				playlistId: item.playlistId,

				loggingContext: item?.loggingContext || null,
				config: { type: item?.musicVideoType },
			});
			return;
		}
	}
	if (item?.endpoint?.pageType.includes("ARTIST")) {
		goto(`/artist/${item?.endpoint?.browseId}`);
	}

	if (
		!isBrowseEndpoint &&
		item.videoId !== undefined &&
		!item?.endpoint?.pageType.includes("ARTIST")
	) {
		await list.initAutoMixSession({
			loggingContext: item?.loggingContext,
			videoId: item.videoId,
			playlistId: item.playlistId,

			config: { type: item?.musicVideoType },
		});
		return;
	} else {
        console.log("return")
		return;
	}
}
