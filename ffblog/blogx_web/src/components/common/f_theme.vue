<script setup lang="ts">
import {theme, setTheme} from "@/components/common/f_theme";

interface Props {
  simple?: boolean
}

const props = defineProps<Props>()

</script>

<template>
  <template v-if="theme === ''">
    <i title="白天模式">
      <icon-sun-fill @click="setTheme('dark')"/>
    </i>
    <span v-if="!props.simple" @click="setTheme('dark')">白天模式</span>
  </template>
  <template v-else>
    <i title="黑夜模式">
      <icon-moon-fill @click="setTheme('')"/>
    </i>
    <span v-if="!props.simple" @click="setTheme('')">黑夜模式</span>
  </template>
</template>

<style scoped>

</style>