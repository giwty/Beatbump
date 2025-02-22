<svelte:options immutable={true} />

<script
	context="module"
	lang="ts"
>
    import {APIClient} from "$lib/api";

    const RE_ALBUM_PLAYLIST_SINGLE = /PLAYLIST|ALBUM|SINGLE/;
	const RE_THUMBNAIL_DIM = /=w\d+-h\d+-/gm;

	const imageErrorHandler = (event: Event) => {
		if (!browser) return;
		const target = event.currentTarget as HTMLImageElement;
		target.onerror = null;

		target.src =
			"data:image/svg+xml;base64,PHN2ZyB4bWxucz0iaHR0cDovL3d3dy53My5vcmcvMjAwMC9zdmciIHN0eWxlPSJpc29sYXRpb246aXNvbGF0ZSIgdmlld0JveD0iMCAwIDI1NiAyNTYiIHdpZHRoPSIyNTZwdCIgaGVpZ2h0PSIyNTZwdCI+PGRlZnM+PGNsaXBQYXRoIGlkPSJwcmVmaXhfX2EiPjxwYXRoIGQ9Ik0wIDBoMjU2djI1NkgweiIvPjwvY2xpcFBhdGg+PC9kZWZzPjxnIGNsaXAtcGF0aD0idXJsKCNwcmVmaXhfX2EpIj48cGF0aCBmaWxsPSIjNDI0MjQyIiBkPSJNMCAwaDI1NnYyNTZIMHoiLz48ZyBjbGlwLXBhdGg9InVybCgjcHJlZml4X19iKSI+PHRleHQgdHJhbnNmb3JtPSJ0cmFuc2xhdGUoMTA1LjU0IDE2Ni43OTQpIiBmb250LWZhbWlseT0ic3lzdGVtLXVpLC1hcHBsZS1zeXN0ZW0sQmxpbmtNYWNTeXN0ZW1Gb250LCZxdW90O1NlZ29lIFVJJnF1b3Q7LFJvYm90byxPeHlnZW4sVWJ1bnR1LENhbnRhcmVsbCwmcXVvdDtPcGVuIFNhbnMmcXVvdDssJnF1b3Q7SGVsdmV0aWNhIE5ldWUmcXVvdDssc2Fucy1zZXJpZiIgZm9udC13ZWlnaHQ9IjQwMCIgZm9udC1zaXplPSIxMDAiIGZpbGw9IiNmYWZhZmEiPj88L3RleHQ+PC9nPjxkZWZzPjxjbGlwUGF0aCBpZD0icHJlZml4X19iIj48cGF0aCB0cmFuc2Zvcm09InRyYW5zbGF0ZSg5MiA1NC44MzkpIiBkPSJNMCAwaDcydjE0Ni4zMjNIMHoiLz48L2NsaXBQYXRoPjwvZGVmcz48L2c+PC9zdmc+";
	};

	function handleContextMenu(event: MouseEvent, dropdownItems: Dropdown) {
		event.preventDefault();
		window.dispatchEvent(
			new CustomEvent("contextmenu", { detail: "carouselItem" }),
		);

		PopperStore.set({
			items: dropdownItems,
			x: event.pageX,
			y: event.pageY,
			direction: "normal",
		});
	}

	const FILTER_ARTIST_ON_ARTIST_PAGE: ReadonlyArray<string> = [
		"Favorite",
		"Add to Queue",
		"View Artist",
	] as const;
	const FILTER_ALBUM_PLAYLIST_ITEMS: ReadonlyArray<string> = [
		"Favorite",
		"Play Next",
		"View Artist",
	] as const;

	const MENU_HANDLERS = {
		artist: async (ctx: BuildMenuParams) => {
			const { item } = ctx;
			try {
				const artistId = item.artistInfo
					? item.artistInfo?.artist?.[0].browseId
					: item.subtitle[0].browseId;
				if (!artistId)
					throw new Error(
						`Expected a valid artistId string, received ${artistId}`,
					);
				goto(`/artist/${artistId}`);
				await tick();
				window.scrollTo({
					behavior: "smooth",
					top: 0,
					left: 0,
				});
			} catch (e) {
				notify(`Error: ${e}`, "error");
			}
		},
		addToQueue: (ctx: BuildMenuParams) => {
			const { item } = ctx;
			list.setTrackWillPlayNext(item, list.$.value.mix.length);
			notify(`${item.title} has been added to your queue!`, "success");
		},
		playNext: (ctx: BuildMenuParams) => {
			const { item } = ctx;
			list.setTrackWillPlayNext(item, list.position);
			notify(`${item.title} will play next!`, "success");
		},
		startGroupSession: () => showGroupSessionCreator.set(true),
		shareGroupSession: async (ctx: BuildMenuParams) => {
			if (!browser) return;
			const { SITE_ORIGIN_URL } = ctx;
			const shareData = {
				title: `Join ${groupSession.client.displayName}'s Beatbump Session`,

				url: `${SITE_ORIGIN_URL}/session?token=${IsoBase64.toBase64(
					JSON.stringify({
						clientId: groupSession.client.clientId,
						displayName: groupSession.client.displayName,
					}),
				)}`,
			};
			try {
				if (!navigator.canShare) {
					await navigator.clipboard.writeText(shareData.url);
					notify("Link copied successfully", "success");
				} else {
					const share = await navigator.share(shareData);
					notify("Shared successfully", "success");
				}
			} catch (error) {
				notify("Error: " + error, "error");
			}
		},
		addToPlaylist: async (ctx: BuildMenuParams) => {
			// console.log(ctx);
			const { item } = ctx;
			if (item.endpoint?.pageType.match(RE_ALBUM_PLAYLIST_SINGLE)) {
				const response = await APIClient.fetch(
                    `/api/v1/get_queue.json?playlistId=` + item.playlistId,
				);
				const data = await response.json();
				const items: Item[] = data;
				showAddToPlaylistPopper.set({ state: true, item: [...items] });
			} else {
				showAddToPlaylistPopper.set({ state: true, item: item });
			}
		},
		favorite: (ctx: BuildMenuParams) => {
			const { item } = ctx;
			IDBService.sendMessage("create", "favorite", item);
		},
		share: async (ctx: BuildMenuParams) => {
			const { SITE_ORIGIN_URL: $SITE_ORIGIN_URL, item } = ctx;
			const shareData = createShare({
				origin: $SITE_ORIGIN_URL,
				id: item.endpoint?.browseId ?? item.videoId,
				type: item.endpoint?.pageType as SharePageType,
				title: item.title,
			});
			try {
				if (!navigator.canShare) {
					await navigator.clipboard.writeText(shareData.url);
					notify("Link copied successfully", "success");
				} else {
					await navigator.share(shareData);
					notify("Shared successfully", "success");
				}
			} catch (error) {
				notify("Failed to share: " + error, "error");
			}
		},
	};
	const buildMenu = (ctx: BuildMenuParams) =>
		buildDropdown()
			.add("View Artist", MENU_HANDLERS.artist.bind(MENU_HANDLERS.artist, ctx))
			.add(
				"Add to Queue",
				MENU_HANDLERS.addToQueue.bind(MENU_HANDLERS.addToQueue, ctx),
			)
			.add(
				"Play Next",
				MENU_HANDLERS.playNext.bind(MENU_HANDLERS.playNext, ctx),
			)
			.add(
				"Add to Playlist",
				MENU_HANDLERS.addToPlaylist.bind(MENU_HANDLERS.addToPlaylist, ctx),
			)
			.add("Favorite", MENU_HANDLERS.favorite.bind(MENU_HANDLERS.favorite, ctx))
			.add(
				"Start Group Session",
				MENU_HANDLERS.startGroupSession.bind(
					MENU_HANDLERS.startGroupSession,
					ctx,
				),
			)
			.add("Share", MENU_HANDLERS.share.bind(MENU_HANDLERS.share, ctx))
			.build();
</script>

<script lang="ts">
	import { goto } from "$app/navigation";
	import Loading from "$components/Loading/Loading.svelte";
	// import { groupSession } from "$lib/stores";
	import { IDBService } from "$lib/workers/db/service";

	import { browser } from "$app/environment";
	import { buildDropdown, type Dropdown } from "$lib/configs/dropdowns.config";
	import { APIParams } from "$lib/constants";
	import { createShare, type SharePageType } from "$lib/shared/createShare";
	import list from "$lib/stores/list";
	import type { Item, Thumbnail } from "$lib/types";
	import type { BuildMenuParams } from "$lib/types/common";
	import type { IListItemRenderer } from "$lib/types/musicListItemRenderer";
	import { IsoBase64, noop, notify } from "$lib/utils";
	import { groupSession } from "$stores/sessions";
	import {
		showAddToPlaylistPopper,
		showGroupSessionCreator,
	} from "$stores/stores";
	import { SITE_ORIGIN_URL } from "$stores/url";
	import { tick } from "svelte";
	import { PopperButton, PopperStore } from "../Popper";
	import { clickHandler } from "./functions";

	export let index: number;
	export let item: IListItemRenderer;
	export let type = "";
	export let kind = "";
	export let aspectRatio: string;
	export let isBrowseEndpoint = false;
	export let nofollow = false;

	const { hasActiveSessionState } = groupSession;
	let loading = false;

	$: RATIO_RECT =
		(aspectRatio?.includes("TWO_LINE_STACK") &&
			kind !== "Fans might also like") ||
		aspectRatio?.includes("16_9")
			? true
			: false;
	$: ASPECT_RATIO = !RATIO_RECT ? "1x1" : "16x9";

	$: ctx = {
		item,
		SITE_ORIGIN_URL: $SITE_ORIGIN_URL,
		dispatch: noop,
		idx: index,
		page: item.endpoint?.pageType as Exclude<SharePageType, null>,
	};

	$: DropdownItems = buildMenu(ctx);
	$: {
		if (
			type === "artist" ||
			(item.endpoint &&
				item.endpoint.pageType?.includes("MUSIC_PAGE_TYPE_ARTIST"))
		) {
			DropdownItems = DropdownItems.filter((item) =>
				FILTER_ARTIST_ON_ARTIST_PAGE.includes(item.text),
			);
		}
		if (item.endpoint?.pageType) {
			DropdownItems = item?.endpoint?.pageType.match(RE_ALBUM_PLAYLIST_SINGLE)
				? [
						{
							action: () => {
								list.initPlaylistSession({
									playlistId: item.playlistId,
									params: APIParams.finite,
									index: 0,
								});
							},
							icon: "shuffle",
							text: "Shuffle",
						},
						{
							action: () => {
								list.setTrackWillPlayNext(item, $list.position);
							},
							icon: "queue",
							text: "Play Next",
						},
						{
							action: () => {
								list.initPlaylistSession({
									playlistId: "RDAMPL" + item.playlistId,
									params: APIParams.finite,
									index: 0,
								});
							},
							icon: "radio",
							text: "Start Radio",
						},
						...DropdownItems.filter(
							(item) => !FILTER_ALBUM_PLAYLIST_ITEMS.includes(item.text),
						),
				  ]
				: DropdownItems.filter((item) =>
						FILTER_ALBUM_PLAYLIST_ITEMS.includes(item.text),
				  );
		}
	}
	$: {
		if (Array.isArray(DropdownItems)) {
			if ($hasActiveSessionState === true) {
				const idxOfSessionItem = DropdownItems.findIndex((item) =>
					item.text.includes("Group Session"),
				);

				DropdownItems[idxOfSessionItem] = {
					text: "Share Group Session",
					action: MENU_HANDLERS.shareGroupSession.bind(
						MENU_HANDLERS.shareGroupSession,
						ctx,
					),
					icon: "share",
				};
			} else {
				const idxOfSessionItem = DropdownItems.findIndex((item) =>
					item.text.includes("Group Session"),
				);

				DropdownItems[idxOfSessionItem] = {
					text: "Start Group Session",
					action: MENU_HANDLERS.startGroupSession.bind(
						MENU_HANDLERS.startGroupSession,
						ctx,
					),
					icon: "users",
				};
			}
		}
		// eslint-disable-next-line no-self-assign
		DropdownItems = DropdownItems;
	}
	$: srcImg = Array.isArray(item?.thumbnails)
		? (item?.thumbnails.at(0) as Thumbnail)
		: { width: 0, height: 0, url: "", placeholder: "" };

	$: srcImg.url =
		srcImg.width < 100
			? srcImg.url.replace(RE_THUMBNAIL_DIM, "=w240-h240-")
			: srcImg.url;

	$: isArtistKind = kind === "Fans might also like";
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
<article
	class="item item{ASPECT_RATIO}"
	on:contextmenu={(event) => handleContextMenu(event, DropdownItems)}
	on:click|stopPropagation={async () => {
		loading = true;
		loading = await clickHandler({ isBrowseEndpoint, index, item, kind, type });
	}}
>
	<section class="item-thumbnail-wrapper img{ASPECT_RATIO}">
		<div
			class="item-thumbnail img{ASPECT_RATIO}"
			class:isArtistKind
		>
			<!-- svelte-ignore a11y-no-noninteractive-tabindex -->
			<div
				class="image img{ASPECT_RATIO}"
				class:isArtistKind
				tabindex="0"
				title={item.title}
			>
				{#if loading}
					<Loading />
				{/if}
				<img
					alt="thumbnail img{ASPECT_RATIO}"
					on:error={(e) => imageErrorHandler(e)}
					loading={index >= 3 ? "lazy" : "eager"}
					decoding="async"
					width={srcImg.width}
					height={srcImg.height}
					src={index >= 3 ? srcImg.placeholder : srcImg.url}
					data-src={index >= 3 ? srcImg.url : null}
				/>
			</div>
			<div class="item-menu">
				<PopperButton
					tabindex={0}
					bind:items={DropdownItems}
				/>
			</div>
		</div>
	</section>
	<div
		class="item-title"
		class:isArtistKind
	>
		<span class="h1 link">
			{item.title}
		</span>
		{#if item.subtitle}
			<div class="subtitles secondary">
				{#each item.subtitle as sub}
					{#if !sub?.browseId}
						<span>{sub.text}</span>
					{:else}
						<a
							on:click|stopPropagation|preventDefault={() => {
								goto("/artist/" + sub?.browseId);
							}}
							rel={nofollow ? "nofollow" : ""}
							href={"/artist/" + sub?.browseId}><span>{sub.text}</span></a
						>
					{/if}
				{/each}
			</div>
		{/if}
	</div>
</article>

<style lang="scss">
	@import "../../../global/redesign/utility/mixins/media-query";
	@import "../../../global/redesign/utility/mixins/old";

	article {
		--thumbnail-radius: clamp(
			4pt,
			calc(var(--column-width, 0px) - 256px),
			16px
		);

		padding: 0.75em;
		margin-bottom: 1em;
		// min-width: 100%;

		scroll-snap-align: start;
		width: var(--column-width);
		contain: layout paint style;
		flex: 0 1;
		display: flex;
		flex-direction: column;
		// flex: 1 0;
		@media (hover: hover) {
			&:hover {
				> :where(.image)::before {
					background-image: linear-gradient(
						rgb(0 0 0 / 88%) 0%,
						rgb(0 0 0 / 72%) 8%,
						rgb(0 0 0 /42.5%) 19%,
						rgb(0 0 0 / 31%) 24%,
						rgb(0 0 0 / 28%) 32%,
						rgb(0 0 0 / 25%) 51%,
						rgb(0 0 0 / 20%) 59%,
						rgb(0 0 0 / 14%) 65%,
						rgb(0 0 0 / 10%) 71%,
						rgb(0 0 0 / 8%) 79%,
						rgb(0 0 0 /5%) 83%,
						rgb(0 0 0 / 3%) 92%,
						rgb(0 0 0 / 0%) 100%
					);
					opacity: 0.7;
					z-index: 1;
				}
			}
		}
	}

	.item-thumbnail-wrapper {
		position: relative;
		overflow: hidden;
		display: block;
		margin-bottom: 0.5rem;

		&.img16x9 {
			aspect-ratio: 16/9;
		}

		&.img1x1 {
			aspect-ratio: 1/1;
		}
	}

	:where(.item-title.isArtistKind) {
		text-align: center;
	}

	:where(.image.isArtistKind) {
		width: var(--thumbnail-size);
		height: var(--thumbnail-size);
		border-radius: 99999em !important;
	}

	a {
		color: inherit;
		transition: color 100ms linear;

		&:hover {
			color: #eee;
		}
	}

	:where(.item-title) {
		display: block;
	}

	:where(.link) {
		display: block;
		display: box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
		text-overflow: ellipsis;
		margin-bottom: 0.325em;
	}

	.item1x1 {
		position: relative;
	}

	.item16x9 {
		width: 100%;
	}

	.img1x1 {
		// width: 100%;

		aspect-ratio: 1/1 !important;
		height: var(--thumbnail-size);
		width: var(--thumbnail-size);
	}

	.img16x9 {
		min-width: calc(calc(var(--column-width) * 1));
		width: 100%;
		min-height: var(--thumbnail-size);
		aspect-ratio: 16/9 !important;
	}

	.subtitles {
		display: block;
		display: box;
		-webkit-line-clamp: 2;
		-webkit-box-orient: vertical;
		overflow: hidden;
		text-overflow: ellipsis;
		cursor: pointer;
		@include trim(3);
	}

	.h1 {
		font-size: 1em;
		line-height: 1.25;
		font-weight: 400 !important;
		display: inline;
		@include trim(3);
	}

	.image {
		min-height: 100%;
		position: relative;
		cursor: pointer;
		user-select: none;
		border-radius: var(--thumbnail-radius);
		overflow: hidden;
		display: flex;
		align-items: center;
		contain: paint;
		// will-change: opacity;

		&:focus {
			border: none;
		}

		&::before {
			position: absolute;
			content: "";
			inset: 0;
			background-image: linear-gradient(
				rgb(0 0 0 / 88%) 0%,
				rgb(0 0 0 / 72%) 8%,
				rgb(0 0 0 /42.5%) 19%,
				rgb(0 0 0 / 31%) 24%,
				rgb(0 0 0 / 28%) 32%,
				rgb(0 0 0 / 25%) 51%,
				rgb(0 0 0 / 20%) 59%,
				rgb(0 0 0 / 14%) 65%,
				rgb(0 0 0 / 10%) 71%,
				rgb(0 0 0 / 8%) 79%,
				rgb(0 0 0 /5%) 83%,
				rgb(0 0 0 / 3%) 92%,
				rgb(0 0 0 / 0%) 100%
			);
			pointer-events: none;
			transition: background-image linear 0.1s, opacity linear 0.1s;
			opacity: 0.5;
			z-index: 1;
		}

		&:active:hover::before {
			background-image: linear-gradient(
				rgb(0 0 0 / 88%) 0%,
				rgb(0 0 0 / 72%) 8%,
				rgb(0 0 0 /42.5%) 19%,
				rgb(0 0 0 / 31%) 24%,
				rgb(0 0 0 / 28%) 32%,
				rgb(0 0 0 / 25%) 51%,
				rgb(0 0 0 / 20%) 59%,
				rgb(0 0 0 / 14%) 65%,
				rgb(0 0 0 / 10%) 71%,
				rgb(0 0 0 / 8%) 79%,
				rgb(0 0 0 /5%) 83%,
				rgb(0 0 0 / 3%) 92%,
				rgb(0 0 0 / 0%) 100%
			);

			opacity: 1;
		}

		> :where(img) {
			aspect-ratio: inherit;
			user-select: none;
			contain: content;
			width: inherit;
			width: 100%;
			height: 100%;
			object-fit: cover;
		}

		@media screen and (max-width: 640px) {
			&::before {
				background-image: linear-gradient(
					rgb(0 0 0 / 88%) 0%,
					rgb(0 0 0 / 72%) 8%,
					rgb(0 0 0 /42.5%) 19%,
					rgb(0 0 0 / 31%) 24%,
					rgb(0 0 0 / 28%) 32%,
					rgb(0 0 0 / 25%) 51%,
					rgb(0 0 0 / 20%) 59%,
					rgb(0 0 0 / 14%) 65%,
					rgb(0 0 0 / 10%) 71%,
					rgb(0 0 0 / 8%) 79%,
					rgb(0 0 0 /5%) 83%,
					rgb(0 0 0 / 3%) 92%,
					rgb(0 0 0 / 0%) 100%
				);
				opacity: 0.7;
				z-index: 1;
			}
		}
	}

	:where(.item) {
		isolation: isolate;
	}

	:where(.image):hover {
		+ :where(.item-menu) {
			opacity: 1 !important;
		}

		& {
			box-shadow: 0 -1px 27px -16px #000 !important;
		}

		&:active:hover {
			box-shadow: 0 -1px 27px -12px #000 !important;
		}
	}

	.item-menu {
		position: absolute;
		right: 0;
		top: 0;
		z-index: 5;
		isolation: isolate;
		margin: 0.25rem;
		opacity: 0;
		transition: 0.1s opacity linear;

		&:focus-visible,
		&:focus-within,
		&:hover {
			opacity: 1;
		}

		@media screen and (max-width: 640px) {
			opacity: 1;
		}

		@media screen and (hover: none) {
			opacity: 1;
		}
	}

	@mixin active {
		> .image {
			&::before {
				background-image: linear-gradient(
					rgb(0 0 0 / 88%) 0%,
					rgb(0 0 0 / 72%) 8%,
					rgb(0 0 0 /42.5%) 19%,
					rgb(0 0 0 / 31%) 24%,
					rgb(0 0 0 / 28%) 32%,
					rgb(0 0 0 / 25%) 51%,
					rgb(0 0 0 / 20%) 59%,
					rgb(0 0 0 / 14%) 65%,
					rgb(0 0 0 / 10%) 71%,
					rgb(0 0 0 / 8%) 79%,
					rgb(0 0 0 /5%) 83%,
					rgb(0 0 0 / 3%) 92%,
					rgb(0 0 0 / 0%) 100%
				);
				opacity: 0.7;
				z-index: 1;
			}
		}
	}

	.item-thumbnail {
		cursor: pointer;
		contain: paint;

		&:focus-visible,
		&:hover,
		&:focus-within {
			@include active;
		}

		position: absolute;
		top: 0;
		height: 100%;
		overflow: hidden;
		border-radius: var(--thumbnail-radius);
	}

	.hidden {
		display: none !important;
		visibility: hidden !important;
	}

	.image,
	img {
		&:focus {
			outline: none;
		}
	}
</style>
