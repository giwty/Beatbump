<script lang="ts">
	import Header from "$lib/components/Layouts/Header.svelte";
	import { onMount } from "svelte";
	import { browser } from "$app/environment";
	import { APIClient } from "$lib/api";
	import list from "$lib/stores/list";

	let tasks: any[] = [];
	let expandedTaskID: number | null = null;
	let taskTracks: any[] = [];

	const fetchTasks = async () => {
		if (!browser) return;
		const res = await APIClient.fetch("/api/v1/downloads");
		if (res.ok) {
			tasks = (await res.json()) || [];
			// If a task is expanded, refresh its tracks too
			if (expandedTaskID) {
				fetchTaskTracks(expandedTaskID);
			}
		}
	};

	const fetchTaskTracks = async (taskID: number) => {
		const res = await APIClient.fetch(`/api/v1/downloads/${taskID}/tracks`);
		if (res.ok) {
			taskTracks = await res.json();
		}
	};

	const toggleTask = async (taskID: number) => {
		if (expandedTaskID === taskID) {
			expandedTaskID = null;
			taskTracks = [];
		} else {
			expandedTaskID = taskID;
			await fetchTaskTracks(taskID);
		}
	};

	const pauseTask = async (taskId: number) => {
		await APIClient.post(`/api/v1/downloads/${taskId}/pause`, {});
		fetchTasks();
	};

	const playTrack = async (track: any) => {
		await list.initAutoMixSession({
			videoId: track.VideoID,
			keyId: 0,
		});
	};

	const resumeTask = async (taskId: number) => {
		await APIClient.post(`/api/v1/downloads/${taskId}/resume`, {});
		fetchTasks();
	};

	const retryTask = async (taskId: number) => {
		await APIClient.post(`/api/v1/downloads/${taskId}/retry`, {});
		fetchTasks();
	};

	onMount(() => {
		fetchTasks();
	});
</script>

<Header
	title="Downloads"
	url="/downloads"
	desc="View download status"
/>

<main class="resp-content-width">
	<div class="controls">
		<button
			on:click={fetchTasks}
			class="refresh-btn">Refresh</button
		>
	</div>
	{#if tasks.length === 0}
		<p>No downloads found.</p>
	{:else}
		<div class="tasks">
			{#each tasks as task}
				<div class="task-container">
					<!-- Parent Task -->
					<div
						class="task parent"
						on:click={() => toggleTask(task.ID)}
					>
						<div class="info">
							<span class="title"
								>{task.PlaylistName || "Unknown Playlist"}</span
							>
							<span class="meta">ID: {task.ReferenceID}</span>
							<span class="meta date"
								>{new Date(task.CreatedAt).toLocaleString()}</span
							>
						</div>

						<div
							class="task-actions"
							on:click|stopPropagation
						>
							{#if task.Status === "processing" || task.Status === "pending"}
								<button
									class="action-btn"
									on:click={() => pauseTask(task.ID)}
								>
									Pause
								</button>
							{/if}
							{#if task.Status === "paused"}
								<button
									class="action-btn"
									on:click={() => resumeTask(task.ID)}
								>
									Resume
								</button>
							{/if}
							{#if task.Status === "completed" && task.Failed > 0}
								<button
									class="action-btn"
									on:click={() => retryTask(task.ID)}
								>
									Retry ({task.Failed} failed)
								</button>
							{/if}
						</div>

						<div class="stats">
							<div class="stat-item">
								<span class="label">Total</span>
								<span class="value">{task.TotalTracks}</span>
							</div>
							<div class="stat-item">
								<span class="label">Done</span>
								<span class="value success">{task.Processed}</span>
							</div>
							<div class="stat-item">
								<span class="label">Failed</span>
								<span class="value error">{task.Failed}</span>
							</div>
						</div>

						<div
							class="status"
							class:completed={task.Status === "completed"}
							class:failed={task.Status === "failed"}
							class:processing={task.Status === "processing" ||
								task.Status === "pending"}
						>
							{task.Status}
						</div>
					</div>

					<!-- Child Tasks (Tracks) -->
					{#if expandedTaskID === task.ID}
						<div class="children">
							{#if taskTracks.length === 0}
								<p class="no-tracks">No tracks found or loading...</p>
							{:else}
								{#each taskTracks as track}
									<div class="child-task">
										{#if track.ThumbnailURL}
											<img
												src={track.ThumbnailURL}
												alt={track.Title}
												class="thumbnail"
											/>
										{/if}
										<div class="info">
											<span class="title">{track.Title}</span>
											<span class="artist">{track.Artist}</span>
										</div>
										<div>
											<button
												class="action-btn"
												on:click={() => playTrack(track)}>Play</button
											>
										</div>
										<div
											class="status small"
											class:completed={track.Status === "completed"}
											class:failed={track.Status === "failed"}
											class:processing={track.Status === "in_progress" ||
												track.Status === "not_started"}
										>
											{track.Status}
										</div>
									</div>
								{/each}
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>
	{/if}
</main>

<style lang="scss">
	.tasks {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}

	.controls {
		display: flex;
		justify-content: flex-end;
		margin-bottom: 1rem;
	}

	.refresh-btn {
		padding: 0.5rem 1rem;
		background: #333;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-weight: bold;

		&:hover {
			background: #444;
		}
	}

	.task-container {
		display: flex;
		flex-direction: column;
		background: rgba(255, 255, 255, 0.05);
		border-radius: 8px;
		overflow: hidden;
	}

	.task {
		padding: 1rem;
		display: flex;
		justify-content: space-between;
		align-items: center;
		gap: 1rem;
		cursor: pointer;
		transition: background 0.2s;

		&:hover {
			background: rgba(255, 255, 255, 0.1);
		}

		&.parent {
			border-bottom: 1px solid rgba(255, 255, 255, 0.05);
		}
	}

	.task-actions {
		display: flex;
		gap: 0.5rem;
		align-items: center;
	}

	.action-btn {
		padding: 0.5rem 1rem;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 4px;
		cursor: pointer;
		font-weight: bold;
		font-size: 0.85rem;
		transition: background 0.2s;

		&:hover {
			background: #0056b3;
		}
	}

	.children {
		background: rgba(0, 0, 0, 0.2);
		padding: 1rem;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		max-height: 400px;
		overflow-y: auto;
	}

	.child-task {
		display: flex;
		align-items: center;
		gap: 1rem;
		padding: 0.5rem;
		background: rgba(255, 255, 255, 0.03);
		border-radius: 4px;
	}

	.thumbnail {
		width: 40px;
		height: 40px;
		border-radius: 4px;
		object-fit: cover;
	}

	.info {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
		flex: 1;
	}

	.title {
		font-weight: bold;
		color: white;
	}

	.artist {
		color: #aaa;
		font-size: 0.9rem;
	}

	.meta {
		color: #666;
		font-size: 0.8rem;

		&.date {
			font-size: 0.7rem;
		}
	}

	.stats {
		display: flex;
		gap: 1rem;
		margin-right: 1rem;
	}

	.stat-item {
		display: flex;
		flex-direction: column;
		align-items: center;

		.label {
			font-size: 0.7rem;
			color: #888;
			text-transform: uppercase;
		}

		.value {
			font-weight: bold;
			font-size: 1.1rem;

			&.success {
				color: #00cd6a;
			}
			&.error {
				color: #ff4d4d;
			}
		}
	}

	.status {
		padding: 0.25rem 0.75rem;
		border-radius: 4px;
		background: #333;
		text-transform: capitalize;
		white-space: nowrap;

		&.completed {
			background: #00cd6a;
			color: black;
		}

		&.failed {
			background: #ff4d4d;
			color: white;
		}

		&.processing {
			background: #007bff;
			color: white;
		}

		&.small {
			font-size: 0.8rem;
			padding: 0.15rem 0.5rem;
		}
	}

	.no-tracks {
		color: #888;
		font-style: italic;
		text-align: center;
		padding: 1rem;
	}
</style>
