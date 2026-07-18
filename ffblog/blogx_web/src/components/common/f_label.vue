<script setup lang="ts">
import type {optionsColorType} from "@/options/options";
import {computed} from "vue";

interface Props {
  value: number | string
  options: optionsColorType[]
  type?: "tag"
}

const props = defineProps<Props>()
const item = computed(() => {
  const obj = props.options.find((value) => value.value === props.value)
  if (obj) {
    return obj
  }
  return {label: "-", color: ""}
})

const {type="tag"} = props


</script>

<template>
  <a-tag v-if="type === 'tag'" :color="item.color">{{ item.label }}</a-tag>
  <span v-else>{{ item.label }}</span>
</template>

<style lang="less">

</style>