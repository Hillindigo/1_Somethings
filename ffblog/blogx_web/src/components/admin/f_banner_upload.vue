<script setup lang="ts">
import {userStorei} from "@/stores/user_store";
import F_cover_cutter from "@/components/web/f_cover_cutter.vue";

const store = userStorei()
interface Props {
  modelValue: string
  placeholder?: string
  shape?:string
  width?: number
  height?:number
  rate: string
}


const props = defineProps<Props>()

const {shape="circle", width = 60, height = 60} = props
const emits = defineEmits(["update:modelValue"])

function inputHandler(val: string) {
  emits("update:modelValue", val)
}

function ok(data: string) {
  emits("update:modelValue", data)
}

</script>

<template>
  <div class="f_image_upload">
    <a-input :model-value="props.modelValue" :placeholder="props.placeholder" @input="inputHandler"></a-input>
   <f_cover_cutter :rate="props.rate" @ok="ok">
     <a-image :src="props.modelValue" :preview="false" :class="shape"  :width="width" :height="height"></a-image>
   </f_cover_cutter>
  </div>
</template>

<style lang="less">
.f_image_upload {
  width: 100%;
  .arco-input-wrapper{
    display: flex;
    margin-bottom: 10px;
  }
  .arco-image{
    &.circle{
      border-radius: 50%;
    }
  }
}
</style>