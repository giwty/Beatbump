<script lang="ts">
	import ListItem, {
		listItemPageContext,
	} from "$lib/components/ListItem/ListItem.svelte";
	import list from "$lib/stores/list";
	import { isPagePlaying, showAddToPlaylistPopper } from "$lib/stores/stores";
	import List from "../_List.svelte";

	import type { ParsedCarousel } from "$api/models/Carousel";
	import Carousel from "$lib/components/Carousel/Carousel.svelte";
	import Header from "$lib/components/Layouts/Header.svelte";
	import InfoBox from "$lib/components/Layouts/InfoBox.svelte";
	import ListInfoBar from "$lib/components/ListInfoBar";
	import { CTX_ListItem, releasePageContext } from "$lib/contexts";
	import type { IListItemRenderer } from "$lib/types/musicListItemRenderer";
	import { notify } from "$lib/utils";
	import { onMount } from "svelte";
	import { writable } from "svelte/store";
	import type { PageData } from "./$types";
    import {APIClient} from "$lib/api";

	export let data: PageData;

	const {
		tracks = [],
		header = {
			thumbnails: [],
			description: "",
			playlistId: "",
			secondSubtitle: [],
			subtitles: [],
			title: "",
		},
		id = "",
		continuations,
		carouselContinuations,
		visitorData,
		key,
	} = data;

	$: ctoken = continuations?.continuation || null;
	$: itct = continuations?.clickTrackingParams || undefined;
	let width = 640;
	let pageTitle = header?.title || "";
	let description: string;
	let isLoading = false;
	let hasData: boolean | null = false;
	let carousel: ParsedCarousel<"twoRowItem">;

	const trackStore = writable<IListItemRenderer[]>([]);

	trackStore.set(tracks);

	CTX_ListItem.set({
		innerWidth: width,
		parentPlaylistId: id,
		visitorData: data?.visitorData,
	});

	pageTitle =
		pageTitle.length > 64
			? pageTitle.substring(0, 64) + "..."
			: header?.title || "";
	description =
		header?.description !== undefined
			? header?.description.length > 240
				? header?.description.substring(0, 240) + "..."
				: header?.description
			: "";
	$: !import.meta.env.SSR && console.log(data);
	const getCarousel = async () => {
		if (!carouselContinuations) return;
		const response = await APIClient.fetch(
            `/api/v1/playlist.json` +
				"?ref=" +
				id +
				`${
					carouselContinuations
						? `&ctoken=${encodeURIComponent(
								carouselContinuations?.continuation,
						  )}`
						: ""
				}` +
				"&itct=" +
				carouselContinuations?.clickTrackingParams,
		);
		const data = await response.json();

		if (data?.carousel) {
			carousel = { ...data?.carousel };
		}
	};
	const getContinuation = async () => {
		if (isLoading || hasData) return;
		if (!itct || !ctoken) {
			getCarousel();
			hasData = true;
			return;
		}

		try {
			isLoading = true;
			const response = await APIClient.fetch(`/api/v1/playlist.json` +
					"?ref=" +
					id +
					"&visitorData=" +
					visitorData +
					`${
						ctoken
							? `&ctoken=${encodeURIComponent(encodeURIComponent(ctoken))}`
							: ""
					}` +
					"&itct=" +
					itct,
			);
			const data = await response.json();
			const continuationItems = data.tracks;
			// Continuations check
			if (data.continuations) {
				/*
					if response has coninuations object, set the new ITCT and Ctoken
					update tracks
				*/
				ctoken = data.continuations.continuation;
				itct = data.continuations.clickTrackingParams;
				trackStore.update((t) => [...t, ...continuationItems]);
				isLoading = false;
				hasData = data.length === 0;
				return hasData;
			} else {
				/*
					if no continuations object is found, set:
					- ctoken to null
					- itct to defined
				*/
				ctoken = null;
				itct = undefined;
				getCarousel();
				hasData = null;
				isLoading = false;
				trackStore.update((t) =>
					[...t, ...continuationItems].filter((item) => {
						if (item !== null || item !== undefined) {
							return item;
						}
					}),
				);
			}
			return !isLoading;
		} catch (error) {
			hasData = null;
			isLoading = false;
		}
	};
	const setId = () => isPagePlaying.add(id);
	let value;
	const options = [
		{
			label: "Unsorted",
			params: "nosort",
			// eslint-disable-next-line @typescript-eslint/no-empty-function
			action: () => {},
		},
		{
			label: "Artist (A-Z)",
			params: "a-az",
			action: () => {
				$trackStore = [
					...$trackStore.sort((a, b) => {
						const itemA = a.artistInfo.artist?.[0]?.text?.toLowerCase() || "";
						const itemB = b.artistInfo.artist?.[0]?.text?.toLowerCase() || "";
						if (itemA < itemB) {
							return -1;
						}
						if (itemA > itemB) {
							return 1;
						}
						return 0;
					}),
				];
			},
		},
		{
			label: "Artist (Z-A)",
			params: "a-za",
			action: () => {
				$trackStore = [
					...$trackStore.sort((a, b) => {
						const itemA = a.artistInfo.artist?.[0]?.text?.toLowerCase() || "";
						const itemB = b.artistInfo.artist?.[0]?.text?.toLowerCase() || "";
						if (itemA < itemB) {
							return 1;
						}
						if (itemA > itemB) {
							return -1;
						}
						return 0;
					}),
				];
			},
		},
		{
			label: "Title (A-Z)",
			params: "t-az",
			action: () => {
				$trackStore = [
					...$trackStore.sort((a, b) => {
						const itemA = a.title.toLowerCase();
						const itemB = b.title.toLowerCase();
						if (itemA < itemB) {
							return -1;
						}
						if (itemA > itemB) {
							return 1;
						}
						return 0;
					}),
				];
			},
		},
		{
			label: "Title (Z-A)",
			params: "t-za",
			action: () => {
				$trackStore = [
					...$trackStore.sort((a, b) => {
						const itemA = a.title.toLowerCase();
						const itemB = b.title.toLowerCase();
						if (itemA < itemB) {
							return 1;
						}
						if (itemA > itemB) {
							return -1;
						}
						return 0;
					}),
				];
			},
		},
	];
	let filter = value ? value : 0;

	onMount(() => {
		const cb = listItemPageContext.add("playlist");
		return () => {
			cb();
		};
	});

	releasePageContext.set({ page: "playlist" });
</script>

<svelte:window bind:innerWidth={width} />
{#if header.title !== "error"}
	<Header
		title={header?.title}
		url={`${key}`}
		desc={description}
		image={header?.thumbnails !== null
			? header?.thumbnails?.[header?.thumbnails?.length - 1]?.url
			: undefined}
	/>
{/if}
<main>
	{#if header.title !== "error"}
		<InfoBox
			subtitles={header?.subtitles}
			secondSubtitle={header?.secondSubtitle}
			thumbnail={header?.thumbnails !== null
				? header?.thumbnails?.[header?.thumbnails?.length - 1].url.replace(
						/=(w(\d+))-(h(\d+))/g,
						"=w512-h512",
				  )
				: undefined}
			title={pageTitle}
			{description}
			buttons={[
				{
					action: () => {
						setId();
						// list.startPlaylist(header.playlistId)
						list.initAutoMixSession({
							playlistId: header.playlistId,
							config: { playerParams: "wAEB8gECGAE%3D" },
						});
					},
					icon: "shuffle",
					text: "Shuffle",
				},
				{
					action: () => {
						setId();
						list.initPlaylistSession({
							playlistId: "RDAMPL" + header.playlistId,
							index: 0,
							params: "wAEB8gECGAE%3D",
						});
					},
					icon: "play",
					type: "outlined",
					text: "Start Radio",
				},
				{
					// eslint-disable-next-line @typescript-eslint/no-empty-function
					action: () => {},
					icon: { name: "dots", size: "1.25rem" },
					text: "",
					type: "icon",
				},
			]}
			on:addqueue={() => {
				setId();
				list.initPlaylistSession({ playlistId: header.playlistId, index: 0 });

				notify(`${pageTitle} added to queue!`, "success");
			}}
			on:playlistAdd={async () => {
				const response = await APIClient.fetch(
					`/api/v1/get_queue.json?playlistId=` + header?.playlistId,
				);
				const data = await response.json();
				const items = data;
				showAddToPlaylistPopper.set({ state: true, item: [...items] });
			}}
		/>
	{/if}
	<ListInfoBar
		bind:value={filter}
		on:change={async () => {
			options[filter].action();
		}}
		{options}
	/>

	<List
		bind:items={$trackStore}
		bind:isLoading
		on:getMore={async () => await getContinuation()}
		bind:hasData
		let:item
		let:index
	>
		<ListItem
			on:setPageIsPlaying={() => {
				setId();
			}}
			{item}
			idx={index}
		/>
	</List>
	<footer>
		{#if carousel}
			<Carousel
				header={carousel?.header}
				items={carousel?.items}
				type="home"
				isBrowseEndpoint={false}
				nofollow
			/>
		{/if}
	</footer>
</main>
