import { useState } from "react";
import {
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
  StatusBar,
  KeyboardAvoidingView,
  Platform,
  ActivityIndicator,
  Dimensions,
} from "react-native";
import { router } from "expo-router";

// ─── Types ────────────────────────────────────────────────────────────────────

type Role = "sorter" | "courier";

interface RoleOption {
  id: Role;
  label: string;
  sublabel: string;
  icon: string;
}

// ─── Constants ────────────────────────────────────────────────────────────────

const ROLES: RoleOption[] = [
  {
    id: "sorter",
    label: "Penyortir",
    sublabel: "Kelola Paket",
    icon: "⊞",
  },
  {
    id: "courier",
    label: "Kurir",
    sublabel: "Lihat rute & konfirmasi",
    icon: "◎",
  },
];

const { width } = Dimensions.get("window");

// ─── Mock auth ────────────────────────────────────────────────────────────────
// Ganti dengan call ke FastAPI: POST /auth/login

async function mockLogin(
  email: string,
  password: string,
  role: Role,
): Promise<{ token: string; role: Role; name: string }> {
  await new Promise((r) => setTimeout(r, 1200));

  if (!email || !password) throw new Error("Email dan password wajib diisi.");
  if (password.length < 4) throw new Error("Password minimal 4 karakter.");

  return {
    token: "mock-jwt-" + role,
    role,
    name: role === "sorter" ? "Budi Santoso" : "Eko Prasetyo",
  };
}

// ─── Component ────────────────────────────────────────────────────────────────

export default function LoginScreen() {
  const [selectedRole, setSelectedRole] = useState<Role>("sorter");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [showPassword, setShowPassword] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleLogin() {
    setError("");
    setLoading(true);
    try {
      const result = await mockLogin(email, password, selectedRole);
      // TODO: simpan token ke SecureStore
      // await SecureStore.setItemAsync("token", result.token);

      // Role-based routing
      if (result.role === "sorter") {
        router.replace("/(sorter)/dashboard");
      } else {
        router.replace("/(courier)/dashboard");
      }
    } catch (err: any) {
      setError(err.message ?? "Login gagal. Coba lagi.");
    } finally {
      setLoading(false);
    }
  }

  return (
    <View style={styles.root}>
      <StatusBar barStyle="dark-content" backgroundColor="#F7F6F2" />

      <KeyboardAvoidingView
        style={styles.inner}
        behavior={Platform.OS === "ios" ? "padding" : "height"}
      >
        {/* ── Header ── */}
        <View style={styles.header}>
          <View style={styles.logoMark}>
            <Text style={styles.logoIcon}>⬡</Text>
          </View>
          <Text style={styles.appName}>RouteIQ</Text>
          <Text style={styles.appTagline}>Sistem pengiriman cerdas</Text>
        </View>

        {/* ── Role selector ── */}
        <View style={styles.section}>
          <Text style={styles.sectionLabel}>Masuk sebagai</Text>
          <View style={styles.roleRow}>
            {ROLES.map((role) => {
              const active = selectedRole === role.id;
              return (
                <TouchableOpacity
                  key={role.id}
                  style={[styles.roleCard, active && styles.roleCardActive]}
                  onPress={() => {
                    setSelectedRole(role.id);
                    setError("");
                  }}
                  activeOpacity={0.7}
                  accessibilityRole="radio"
                  accessibilityState={{ checked: active }}
                  accessibilityLabel={`${role.label}, ${role.sublabel}`}
                >
                  <Text
                    style={[styles.roleIcon, active && styles.roleIconActive]}
                  >
                    {role.icon}
                  </Text>
                  <Text
                    style={[styles.roleLabel, active && styles.roleLabelActive]}
                  >
                    {role.label}
                  </Text>
                  <Text
                    style={[
                      styles.roleSublabel,
                      active && styles.roleSublabelActive,
                    ]}
                  >
                    {role.sublabel}
                  </Text>
                  {active && <View style={styles.roleActiveDot} />}
                </TouchableOpacity>
              );
            })}
          </View>
        </View>

        {/* ── Form ── */}
        <View style={styles.section}>
          <Text style={styles.sectionLabel}>Detail akun</Text>

          <View style={styles.inputGroup}>
            <Text style={styles.inputLabel}>Email</Text>
            <TextInput
              style={styles.input}
              value={email}
              onChangeText={(v) => {
                setEmail(v);
                setError("");
              }}
              placeholder="nama@routeiq.id"
              placeholderTextColor="#B4B2A9"
              keyboardType="email-address"
              autoCapitalize="none"
              autoCorrect={false}
              returnKeyType="next"
              editable={!loading}
            />
          </View>

          <View style={[styles.inputGroup, { marginTop: 12 }]}>
            <Text style={styles.inputLabel}>Password</Text>
            <View style={styles.passwordRow}>
              <TextInput
                style={[styles.input, styles.passwordInput]}
                value={password}
                onChangeText={(v) => {
                  setPassword(v);
                  setError("");
                }}
                placeholder="••••••••"
                placeholderTextColor="#B4B2A9"
                secureTextEntry={!showPassword}
                returnKeyType="done"
                onSubmitEditing={handleLogin}
                editable={!loading}
              />
              <TouchableOpacity
                style={styles.eyeBtn}
                onPress={() => setShowPassword((p) => !p)}
                accessibilityLabel={
                  showPassword ? "Sembunyikan password" : "Tampilkan password"
                }
              >
                <Text style={styles.eyeIcon}>{showPassword ? "○" : "●"}</Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>

        {/* ── Error ── */}
        {error !== "" && (
          <View style={styles.errorBox} accessibilityRole="alert">
            <Text style={styles.errorIcon}>!</Text>
            <Text style={styles.errorText}>{error}</Text>
          </View>
        )}

        {/* ── CTA ── */}
        <TouchableOpacity
          style={[styles.loginBtn, loading && styles.loginBtnLoading]}
          onPress={handleLogin}
          activeOpacity={0.85}
          disabled={loading}
          accessibilityRole="button"
          accessibilityLabel="Masuk"
          accessibilityState={{ disabled: loading, busy: loading }}
        >
          {loading ? (
            <ActivityIndicator color="#F7F6F2" size="small" />
          ) : (
            <Text style={styles.loginBtnText}>
              Masuk sebagai {ROLES.find((r) => r.id === selectedRole)?.label}
            </Text>
          )}
        </TouchableOpacity>

        {/* ── Footer ── */}
        <Text style={styles.footer}>
          Akun dibuat oleh admin. Hubungi supervisor jika lupa password.
        </Text>
      </KeyboardAvoidingView>
    </View>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const CARD_WIDTH = (width - 48 - 10) / 2;

const styles = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: "#F7F6F2",
  },
  inner: {
    flex: 1,
    paddingHorizontal: 24,
    paddingTop: 64,
    paddingBottom: 32,
  },

  // Header
  header: {
    alignItems: "center",
    marginBottom: 36,
  },
  logoMark: {
    width: 52,
    height: 52,
    borderRadius: 14,
    backgroundColor: "#1D1D1B",
    alignItems: "center",
    justifyContent: "center",
    marginBottom: 12,
  },
  logoIcon: {
    fontSize: 26,
    color: "#F7F6F2",
    lineHeight: 30,
  },
  appName: {
    fontSize: 24,
    fontWeight: "600",
    color: "#1D1D1B",
    letterSpacing: -0.5,
  },
  appTagline: {
    fontSize: 13,
    color: "#888780",
    marginTop: 2,
  },

  // Section
  section: {
    marginBottom: 20,
  },
  sectionLabel: {
    fontSize: 11,
    fontWeight: "500",
    color: "#888780",
    textTransform: "uppercase",
    letterSpacing: 0.8,
    marginBottom: 10,
  },

  // Role cards
  roleRow: {
    flexDirection: "row",
    gap: 10,
  },
  roleCard: {
    width: CARD_WIDTH,
    backgroundColor: "#FFFFFF",
    borderRadius: 12,
    borderWidth: 1,
    borderColor: "#E0DED6",
    padding: 14,
    alignItems: "flex-start",
    position: "relative",
  },
  roleCardActive: {
    borderColor: "#1D1D1B",
    borderWidth: 1.5,
    backgroundColor: "#FFFFFF",
  },
  roleIcon: {
    fontSize: 22,
    color: "#B4B2A9",
    marginBottom: 8,
  },
  roleIconActive: {
    color: "#1D1D1B",
  },
  roleLabel: {
    fontSize: 15,
    fontWeight: "500",
    color: "#888780",
    marginBottom: 2,
  },
  roleLabelActive: {
    color: "#1D1D1B",
  },
  roleSublabel: {
    fontSize: 11,
    color: "#B4B2A9",
  },
  roleSublabelActive: {
    color: "#5F5E5A",
  },
  roleActiveDot: {
    position: "absolute",
    top: 12,
    right: 12,
    width: 7,
    height: 7,
    borderRadius: 4,
    backgroundColor: "#1D1D1B",
  },

  // Inputs
  inputGroup: {
    gap: 6,
  },
  inputLabel: {
    fontSize: 13,
    fontWeight: "500",
    color: "#444441",
  },
  input: {
    height: 46,
    backgroundColor: "#FFFFFF",
    borderRadius: 10,
    borderWidth: 1,
    borderColor: "#E0DED6",
    paddingHorizontal: 14,
    fontSize: 15,
    color: "#1D1D1B",
  },
  passwordRow: {
    position: "relative",
  },
  passwordInput: {
    paddingRight: 46,
  },
  eyeBtn: {
    position: "absolute",
    right: 0,
    top: 0,
    height: 46,
    width: 46,
    alignItems: "center",
    justifyContent: "center",
  },
  eyeIcon: {
    fontSize: 12,
    color: "#888780",
  },

  // Error
  errorBox: {
    flexDirection: "row",
    alignItems: "center",
    gap: 8,
    backgroundColor: "#FCEBEB",
    borderRadius: 10,
    paddingVertical: 10,
    paddingHorizontal: 14,
    marginBottom: 16,
    borderWidth: 1,
    borderColor: "#F7C1C1",
  },
  errorIcon: {
    fontSize: 14,
    fontWeight: "600",
    color: "#791F1F",
  },
  errorText: {
    fontSize: 13,
    color: "#791F1F",
    flex: 1,
    lineHeight: 18,
  },

  // CTA
  loginBtn: {
    height: 50,
    backgroundColor: "#1D1D1B",
    borderRadius: 12,
    alignItems: "center",
    justifyContent: "center",
    marginBottom: 20,
  },
  loginBtnLoading: {
    backgroundColor: "#444441",
  },
  loginBtnText: {
    fontSize: 15,
    fontWeight: "500",
    color: "#F7F6F2",
    letterSpacing: 0.1,
  },

  // Footer
  footer: {
    fontSize: 12,
    color: "#B4B2A9",
    textAlign: "center",
    lineHeight: 18,
  },
});
