<template>
  <div v-if="candidates.length === 0" class="w-[14rem] text-warning">
    {{ $t("custom-approval.issue-review.no-one-matches-role") }}
  </div>

  <div
    v-else
    class="min-w-[8rem] max-w-[12rem] max-h-[18rem] flex flex-col text-control-light overflow-y-hidden"
  >
    <div class="flex-1 overflow-auto text-sm">
      <div
        v-for="user in candidates"
        :key="user.name"
        class="flex items-center py-1 gap-x-1"
        :class="[user.name === currentUser.name && 'font-bold']"
      >
        <PrincipalAvatar :user="user" size="SMALL" />
        <span class="whitespace-nowrap">{{ user.title }}</span>
        <span
          v-if="currentUser.name === user.name"
          class="inline-flex items-center px-1 py-0.5 rounded-lg text-xs font-semibold bg-green-100 text-green-800"
        >
          {{ $t("custom-approval.issue-review.you") }}
        </span>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import PrincipalAvatar from "@/components/PrincipalAvatar.vue";
import { useCurrentUserV1 } from "@/store";
import { WrappedReviewStep } from "@/types";

const props = defineProps<{
  step: WrappedReviewStep;
}>();

const currentUser = useCurrentUserV1();
const candidates = computed(() => props.step.candidates);
</script>
