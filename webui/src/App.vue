<template>
	<!-- Componente radice dell'app-->
	<router-view/>
</template>

<script>
export default {
	mounted(){
		if (!localStorage.getItem('token')) {
			this.$router.replace("/login")
		} else {
			this.$router.replace("/chat")
		}
	}
}
</script>
