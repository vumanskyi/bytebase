<template>
  <NSelect
    v-bind="$attrs"
    :value="environment"
    :options="options"
    :placeholder="$t('environment.select')"
    :filterable="true"
    :filter="filterByName"
    :render-label="renderLabel"
    class="bb-environment-select"
    @update:value="$emit('update:environment', $event)"
  />
</template>

<script lang="ts" setup>
import { NSelect, SelectOption } from "naive-ui";
import { computed, h, watchEffect } from "vue";
import { useI18n } from "vue-i18n";
import { useEnvironmentV1Store, useProjectV1Store } from "@/store";
import { State } from "@/types/proto/v1/common";
import { Environment } from "@/types/proto/v1/environment_service";
import { EnvironmentV1Name } from "../Model";

interface EnvironmentSelectOption extends SelectOption {
  value: string;
  environment: Environment;
}

const props = withDefaults(
  defineProps<{
    environment?: string | undefined;
    defaultEnvironmentName?: string | undefined;
    includeArchived?: boolean;
    showProductionIcon?: boolean;
    filter?: (environment: Environment, index: number) => boolean;
  }>(),
  {
    environment: undefined,
    defaultEnvironmentName: undefined,
    includeArchived: false,
    showProductionIcon: true,
    filter: () => true,
  }
);

defineEmits<{
  (event: "update:environment", id: string | undefined): void;
}>();

const { t } = useI18n();
const projectV1Store = useProjectV1Store();
const environmentV1Store = useEnvironmentV1Store();

const prepare = () => {
  projectV1Store.fetchProjectList(true /* showDeleted */);
};

const rawEnvironmentList = computed(() => {
  const list = environmentV1Store.getEnvironmentList(true /* showDeleted */);
  return list;
});

const combinedEnvironmentList = computed(() => {
  let list = rawEnvironmentList.value.filter((environment) => {
    if (props.includeArchived) return true;
    if (environment.state === State.ACTIVE) return true;
    // ARCHIVED
    if (environment.uid === props.environment) return true;
    return false;
  });

  if (props.filter) {
    list = list.filter(props.filter);
  }

  return list;
});

const options = computed(() => {
  return combinedEnvironmentList.value.map<EnvironmentSelectOption>(
    (environment) => {
      return {
        environment,
        value: environment.uid,
        label: environment.title,
      };
    }
  );
});

const renderLabel = (option: SelectOption) => {
  const { environment } = option as EnvironmentSelectOption;
  return h(EnvironmentV1Name, {
    environment,
    showIcon: props.showProductionIcon,
    link: false,
    suffix:
      props.defaultEnvironmentName === environment.name
        ? `(${t("common.default")})`
        : "",
  });
};

const filterByName = (pattern: string, option: SelectOption) => {
  const { environment } = option as EnvironmentSelectOption;
  pattern = pattern.toLowerCase();
  return (
    environment.name.toLowerCase().includes(pattern) ||
    environment.title.toLowerCase().includes(pattern)
  );
};

watchEffect(prepare);
</script>
