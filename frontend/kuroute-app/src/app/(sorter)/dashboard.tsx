import { StyleSheet, Text, View } from "react-native";

export default function SorterDashboardScreen() {
  return (
    <View style={styles.container}>
      <Text style={styles.title}>Dashboard Penyortir</Text>
      <Text style={styles.subtitle}>
        Kelola order, rute, dan status pengiriman.
      </Text>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    padding: 24,
    backgroundColor: "#F7F6F2",
  },
  title: {
    fontSize: 24,
    fontWeight: "600",
    color: "#1D1D1B",
    marginBottom: 8,
  },
  subtitle: {
    fontSize: 14,
    color: "#66655F",
    textAlign: "center",
  },
});
