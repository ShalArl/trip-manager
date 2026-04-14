/**
 * Password validation rules:
 * - At least 8 characters
 * - At least one uppercase letter
 * - At least one number
 */

export interface PasswordStrength {
  isValid: boolean;
  score: number; // 0-4 (0: invalid, 1: weak, 2: fair, 3: good, 4: strong)
  errors: string[];
  suggestions: string[];
}

const PASSWORD_RULES = {
  MIN_LENGTH: 8,
  HAS_UPPERCASE: /[A-Z]/,
  HAS_LOWERCASE: /[a-z]/,
  HAS_NUMBER: /\d/,
  HAS_SPECIAL: /[!@#$%^&*()_+\-=\[\]{};:'",.<>?\/]/,
};

export function validatePassword(password: string): PasswordStrength {
  const errors: string[] = [];
  const suggestions: string[] = [];
  let score = 0;

  // Check minimum length
  if (password.length < PASSWORD_RULES.MIN_LENGTH) {
    errors.push(`Mindestens ${PASSWORD_RULES.MIN_LENGTH} Zeichen erforderlich`);
  } else {
    score++;
  }

  // Check for uppercase letter
  if (!PASSWORD_RULES.HAS_UPPERCASE.test(password)) {
    errors.push("Mindestens ein Großbuchstabe erforderlich (A-Z)");
  } else {
    score++;
  }

  // Check for number
  if (!PASSWORD_RULES.HAS_NUMBER.test(password)) {
    errors.push("Mindestens eine Zahl erforderlich (0-9)");
  } else {
    score++;
  }

  // Check for lowercase letter (bonus)
  if (!PASSWORD_RULES.HAS_LOWERCASE.test(password)) {
    suggestions.push("Füge Kleinbuchstaben hinzu für ein stärkeres Passwort");
  } else if (score < 4) {
    score++;
  }

  // Check for special characters (bonus)
  if (PASSWORD_RULES.HAS_SPECIAL.test(password)) {
    score = Math.min(4, score + 1);
    suggestions.push("Passwort ist sehr sicher");
  }

  // Ensure score is at least 0 and at most 4
  score = Math.min(4, Math.max(0, score));

  const isValid = errors.length === 0;

  return {
    isValid,
    score,
    errors,
    suggestions,
  };
}

export function getPasswordStrengthLabel(score: number): string {
  switch (score) {
    case 0:
      return "Zu schwach";
    case 1:
      return "Schwach";
    case 2:
      return "Mittel";
    case 3:
      return "Gut";
    case 4:
      return "Sehr stark";
    default:
      return "Unbekannt";
  }
}

export function getPasswordStrengthColor(score: number): string {
  switch (score) {
    case 0:
      return "text-red-600 dark:text-red-400";
    case 1:
      return "text-orange-600 dark:text-orange-400";
    case 2:
      return "text-yellow-600 dark:text-yellow-400";
    case 3:
      return "text-lime-600 dark:text-lime-400";
    case 4:
      return "text-green-600 dark:text-green-400";
    default:
      return "text-zinc-600 dark:text-zinc-400";
  }
}

export function getPasswordStrengthBarColor(score: number): string {
  switch (score) {
    case 0:
      return "bg-red-500";
    case 1:
      return "bg-orange-500";
    case 2:
      return "bg-yellow-500";
    case 3:
      return "bg-lime-500";
    case 4:
      return "bg-green-500";
    default:
      return "bg-zinc-400";
  }
}

