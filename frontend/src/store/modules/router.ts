import { defineStore } from "pinia";
import { RouteLocationNormalized } from "vue-router";
import { RouterSlug } from "@/types";

export const useRouterStore = defineStore("router", {
  // need not to initialize a state since we store everything into localStorage
  // state: () => ({}),

  getters: {
    backPath: () => () => {
      return localStorage.getItem("ui.backPath") || "/";
    },
  },
  actions: {
    setBackPath(backPath: string) {
      localStorage.setItem("ui.backPath", backPath);
      return backPath;
    },
    routeSlug(currentRoute: RouteLocationNormalized): RouterSlug {
      {
        // /issue/:issueSlug
        // Total 2 elements, 2nd element is the issue slug
        const issueComponents = currentRoute.path.match(
          "/issue/([0-9a-zA-Z_-]+)"
        ) || ["/", undefined];
        if (issueComponents[1]) {
          return {
            issueSlug: issueComponents[1],
          };
        }
      }

      {
        // /sql-editor/sheet/:sheetSlug
        // match this route first
        const sqlEditorComponents = currentRoute.path.match(
          "/sql-editor/sheet/([0-9a-zA-Z_-]+)"
        ) || ["/", undefined];

        if (sqlEditorComponents[1]) {
          return {
            sheetSlug: sqlEditorComponents[1],
          };
        }
      }

      {
        // /sql-editor/:connectionSlug
        const sqlEditorComponents = currentRoute.path.match(
          "/sql-editor/([0-9a-zA-Z_-]+)"
        ) || ["/", undefined];

        if (sqlEditorComponents[1]) {
          return {
            connectionSlug: sqlEditorComponents[1],
          };
        }
      }

      return {};
    },
  },
});
