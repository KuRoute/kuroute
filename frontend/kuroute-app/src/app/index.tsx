import { useEffect } from "react";
import { ActivityIndicator, View } from "react-native";
import { useRouter } from "expo-router";

function useAuth() {
  return {
    isLoggedIn: false,
    role: null as string | null,
    isLoading: false,
  };
}

export default function IndexGatekeeper() {
  const router = useRouter();

  const { isLoggedIn, role, isLoading } = useAuth();

  useEffect(() => {
    if (isLoading) return;

    if (!isLoggedIn) {
      // If belum login, lempar ke halaman login
      router.replace("/(auth)/login");
    } else if (role === "sorter") {
      // If sorter, arahkan ke dasbor sorter
      router.replace("/(sorter)/dashboard");
    } else if (role === "courier") {
      // If courier, arahkan ke dasbor courier
      router.replace("/(courier)/dashboard");
    }
  }, [isLoggedIn, role, isLoading]);

  // Loading state
  return (
    <View style={{ flex: 1, justifyContent: "center", alignItems: "center" }}>
      <ActivityIndicator size="large" color="#0000ff" />
    </View>
  );
}
