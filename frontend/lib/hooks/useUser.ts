import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { getMe, updateMe, changePassword } from "@/lib/api/auth";

// Hook zum Laden des Users
export function useUser() {
  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;

  return useQuery({
    queryKey: ["user"],
    queryFn: getMe,
    enabled: !!token,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Hook zum Aktualisieren des Users
export function useUpdateUser() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: updateMe,
    onSuccess: (data) => {
      // Update cache
      queryClient.setQueryData(["user"], data);
    },
  });
}

// Hook zum Ändern des Passworts
export function useChangePassword() {
  return useMutation({
    mutationFn: changePassword,
  });
}


