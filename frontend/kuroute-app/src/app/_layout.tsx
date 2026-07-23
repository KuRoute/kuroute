import { useEffect, useState } from "react";
import { Stack, router } from "expo-router";
import { View, ActivityIndicator } from "react-native";

async function getStoredSession(): Promise<{
  token: string | null;
  role: "sorter" | "courier" | null;
}> {
  // Simulasi async storage
  await new Promise((r) => setTimeout(r, 300));
  return { token: null, role: null };
}

// Root Layout

export default function RootLayout() {
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    getStoredSession().then(({ token, role }) => {
      if (token && role) {
        // Sudah login > langsung ke screen role
        router.replace(
          role === "sorter" ? "/(sorter)/dashboard" : "/(courier)/dashboard",
        );
      } else {
        // Belum login > ke halaman login
        router.replace("/(auth)/login");
      }
      setChecking(false);
    });
  }, []);

  if (checking) {
    return (
      <View
        style={{
          flex: 1,
          alignItems: "center",
          justifyContent: "center",
          backgroundColor: "#F7F6F2",
        }}
      >
        <ActivityIndicator color="#1D1D1B" />
      </View>
    );
  }

  return (
    <Stack screenOptions={{ headerShown: false }}>
      <Stack.Screen name="index" />
      <Stack.Screen name="(auth)/login" />
      <Stack.Screen name="(courier)/dashboard" />
      <Stack.Screen name="(sorter)/dashboard" />
    </Stack>
  );
}
