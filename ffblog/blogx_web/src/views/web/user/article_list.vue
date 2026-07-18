<script setup lang="ts">
import {useRoute} from "vue-router";
import {userBaseStorei} from "@/stores/user_base_store";
import F_category_list from "@/components/web/article/f_category_list.vue";
import F_article_list from "@/components/web/article/f_article_list.vue";

const baseStore = userBaseStorei()
const route = useRoute()

</script>

<template>
  <div class="user_article_list_view">
    <f_category_list :user-id="Number(route.params.id)" :is-me="baseStore.isMe"></f_category_list>
    <f_article_list></f_article_list>
  </div>
</template>

<style lang="less">
.user_article_list_view {
  display: flex;
  height: 100%;


  .f_category_list_com {
    width: 150px;
  }
}

</style>