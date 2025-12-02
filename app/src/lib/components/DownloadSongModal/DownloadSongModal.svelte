<script lang="ts">
	import { APIClient } from "$lib/api";
	import { showDownloadSongPopper } from "$lib/stores/stores";
	import { createEventDispatcher } from "svelte";
	import Modal from "../Modal";
	import { notify } from "$lib/utils";

	$: isShowing = $showDownloadSongPopper.state;
	$: item = $showDownloadSongPopper.item;

	let limit = 0;
	let hasFocus = false;
	const dispatch = createEventDispatcher();

	async function handleDownload() {
		if (!item) return;

		try {
			const params = new URLSearchParams({
				videoId: item.videoId,
				title: item.title,
				limit: limit.toString(),
				artist: item.artistInfo?.artist?.[0]?.text || "",
				album: item.album?.title || "",
				thumbnailUrl: item.thumbnails?.[item.thumbnails.length - 1]?.url || "",
			});

			const res = await APIClient.fetch(
				`/api/v1/download/song?${params.toString()}`,
			);
			if (res.ok) {
				notify("Download queued successfully", "success");
				showDownloadSongPopper.set({ state: false, item: undefined });
			} else {
				const data = await res.json();
				notify(data.message || "Failed to queue download", "error");
			}
		} catch (e) {
			notify("Error queuing download", "error");
			console.error(e);
		}
	}
</script>

{#if isShowing}
	<Modal
		zIndex={1000}
		on:close={() => {
			showDownloadSongPopper.set({ state: false, item: undefined });
			dispatch("close");
		}}
		bind:hasFocus
	>
		<h1 slot="header">Download Song</h1>
		<div class="content">
			<div class="info">
				<h3>{item?.title}</h3>
				<p>{item?.artistInfo?.artist?.[0]?.text}</p>
			</div>

			<div class="form-group">
				<label for="limit">Related Songs to Download (0-500)</label>
				<input
					type="number"
					id="limit"
					bind:value={limit}
					min="0"
					max="500"
					class="input"
				/>
			</div>

			<div class="actions">
				<button
					class="btn secondary"
					on:click={() =>
						showDownloadSongPopper.set({ state: false, item: undefined })}
					>Cancel</button
				>
				<button
					class="btn primary"
					on:click={handleDownload}>Download</button
				>
			</div>
		</div>
	</Modal>
{/if}

<style lang="scss">
	.content {
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 1rem;

		@media screen and (max-width: 640px) {
			padding: 0.75rem;
		}
	}

	.info {
		h3 {
			margin: 0;
			font-size: 1.2rem;
			font-weight: 600;
		}
		p {
			margin: 0.3rem 0 0;
			color: rgba(255, 255, 255, 0.7);
			font-size: 0.95rem;
		}
	}

	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.6rem;

		label {
			font-weight: 500;
			font-size: 0.95rem;
			text-transform: uppercase;
			letter-spacing: 0.03em;
			color: rgba(255, 255, 255, 0.85);
		}

		.input {
			padding: 0.75rem 1rem;
			border-radius: 8px;
			border: 2px solid rgba(255, 255, 255, 0.1);
			background: rgba(0, 0, 0, 0.3);
			color: #fff;
			font-size: 1rem;
			transition: all 0.2s ease;

			&:focus {
				outline: none;
				border-color: #3b82f6;
				background: rgba(0, 0, 0, 0.4);
			}

			&:hover {
				border-color: rgba(255, 255, 255, 0.2);
			}
		}
	}

	.actions {
		display: flex;
		justify-content: flex-end;
		gap: 0.75rem;
		margin-top: 0.5rem;

		@media screen and (max-width: 640px) {
			flex-direction: column-reverse;
			gap: 0.5rem;
		}

		.btn {
			padding: 0.75rem 1.5rem;
			border-radius: 8px;
			border: none;
			cursor: pointer;
			font-weight: 600;
			font-size: 0.95rem;
			transition: all 0.2s ease;
			min-width: 100px;

			@media screen and (max-width: 640px) {
				width: 100%;
				padding: 0.85rem 1.5rem;
			}

			&:hover {
				transform: translateY(-1px);
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
			}

			&:active {
				transform: translateY(0);
			}

			&.primary {
				background: linear-gradient(135deg, #3b82f6 0%, #2563eb 100%);
				color: white;
				box-shadow: 0 2px 8px rgba(59, 130, 246, 0.3);

				&:hover {
					background: linear-gradient(135deg, #2563eb 0%, #1d4ed8 100%);
					box-shadow: 0 4px 16px rgba(59, 130, 246, 0.4);
				}
			}

			&.secondary {
				background: rgba(255, 255, 255, 0.08);
				color: rgba(255, 255, 255, 0.95);
				border: 2px solid rgba(255, 255, 255, 0.2);
				box-shadow: 0 2px 6px rgba(0, 0, 0, 0.15);

				&:hover {
					background: rgba(255, 255, 255, 0.15);
					border-color: rgba(255, 255, 255, 0.35);
					box-shadow: 0 4px 12px rgba(255, 255, 255, 0.1);
				}
			}
		}
	}
</style>
