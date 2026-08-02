interface ValidationRule {
  value: unknown;
  label: string;
  optional?: boolean;
}

export function useFormValidation() {
  function validate(rules: Record<string, ValidationRule>): string | null {
    for (const [, rule] of Object.entries(rules)) {
      const { value, label, optional } = rule;
      if (optional) continue;
      if (value === null || value === undefined) return label;
      if (typeof value === "string" && !value.trim()) return label;
    }
    return null;
  }

  return { validate };
}
