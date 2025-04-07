<template>
	<div class="emoji-picker" v-if="visible" :style="{ top: `${position.y}px`, left: `${position.x}px` }">
		<div class="emoji-grid">
      <span
		  v-for="emoji in emojis"
		  :key="emoji"
		  class="emoji"
		  @click="$emit('select', emoji)"
	  >
        {{ emoji }}
      </span>
		</div>
	</div>
</template>

<script>
export default {
	name: "EmojiPicker",
	props: {
		visible: Boolean,
		position: {
			type: Object,
			required: true,
		},
	},
	data() {
		return {
			emojis: [
				'😀', '😂', '😅', '😍', '😎', '🤔', '😢', '😡', '👍', '👎',
				'👏', '🙏', '🎉', '❤️', '🔥', '💯', '✅', '❌', '⚡', '🌟'
			],
		};
	},
};
</script>

<style scoped>
.emoji-picker {
	position: absolute;
	background: white;
	border: 1px solid #ccc;
	border-radius: 10px;
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
	z-index: 9999;
	padding: 10px;
}
.emoji-grid {
	display: grid;
	grid-template-columns: repeat(5, 1fr);
	gap: 8px;
}
.emoji {
	font-size: 1.5rem;
	cursor: pointer;
	padding: 4px;
	border-radius: 5px;
	transition: background 0.2s;
	text-align: center;
}
.emoji:hover {
	background: #f0f0f0;
}
</style>
