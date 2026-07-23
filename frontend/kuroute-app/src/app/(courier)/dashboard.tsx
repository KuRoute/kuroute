import React, { useState, useRef } from "react";
import {
  View,
  Text,
  TouchableOpacity,
  ScrollView,
  StyleSheet,
  StatusBar,
  Linking,
  Dimensions,
  Platform,
} from "react-native";
import Mapbox, {
  MapView,
  Camera,
  MarkerView,
  ShapeSource,
  LineLayer,
  CircleLayer,
} from "@rnmapbox/maps";

// ─── Mapbox token ─────────────────────────────────────────────────────────────

const MAPBOX_TOKEN = process.env.EXPO_PUBLIC_MAPBOX_ACCESS_TOKEN || "";
Mapbox.setAccessToken(MAPBOX_TOKEN);

// ─── Types ────────────────────────────────────────────────────────────────────

type ScreenState = "idle" | "pick_trip" | "active";

interface Trip {
  id: string;
  number: number;
  recipientName: string;
  address: string;
  distance: string;
  duration: string;
  price: string;
  coordinate: [number, number]; // [lng, lat] — format GeoJSON Mapbox
}

interface DeliveryStep {
  label: string;
  done: boolean;
  active: boolean;
}

// ─── Mock data ────────────────────────────────────────────────────────────────

const DRIVER_NAME = "Budiono Siregar";
const TOTAL_PACKAGES = 102;

// Mapbox pakai [longitude, latitude] — kebalikan dari react-native-maps
const DEPOT: [number, number] = [115.2126, -8.6726];

const TRIPS: Trip[] = [
  {
    id: "t1",
    number: 1,
    recipientName: "Asep Knalpot",
    address: "Pekalongan V Blok 2B No 10",
    distance: "2.4 km",
    duration: "12 menit",
    price: "IDR 2.400",
    coordinate: [115.218, -8.668],
  },
  {
    id: "t2",
    number: 2,
    recipientName: "Dewi Rahayu",
    address: "Jl. Garuda No 45A",
    distance: "1.8 km",
    duration: "9 menit",
    price: "IDR 1.800",
    coordinate: [115.214, -8.664],
  },
  {
    id: "t3",
    number: 3,
    recipientName: "Hendra Wijaya",
    address: "Jl. Raya Kerambitan No 7",
    distance: "3.1 km",
    duration: "15 menit",
    price: "IDR 3.100",
    coordinate: [115.21, -8.66],
  },
];

// GeoJSON LineString untuk polyline rute
const ROUTE_GEOJSON: GeoJSON.FeatureCollection = {
  type: "FeatureCollection",
  features: [
    {
      type: "Feature",
      properties: {},
      geometry: {
        type: "LineString",
        coordinates: [
          DEPOT,
          [115.215, -8.671],
          [115.217, -8.669],
          TRIPS[0].coordinate,
        ],
      },
    },
  ],
};

// GeoJSON stop markers untuk CircleLayer
const STOPS_GEOJSON: GeoJSON.FeatureCollection = {
  type: "FeatureCollection",
  features: TRIPS.map((trip, i) => ({
    type: "Feature" as const,
    properties: { number: i + 1, id: trip.id },
    geometry: {
      type: "Point" as const,
      coordinates: trip.coordinate,
    },
  })),
};

const STEPS: DeliveryStep[] = [
  { label: "Pesanan dalam pengantaran", done: true, active: false },
  { label: "Kurir sedang menuju alamatmu", done: false, active: true },
  { label: "Kurir telah mengantar pesananmu", done: false, active: false },
];

// ─── Constants ────────────────────────────────────────────────────────────────

const BLUE = "#0D90D4";
const BLUE_DARK = "#0A7AB8";
const BLUE_LIGHT = "#E8F5FC";
const WHITE = "#FFFFFF";
const TEXT_DARK = "#1A1A1A";
const TEXT_MID = "#555555";
const TEXT_LIGHT = "#999999";
const { width: SW } = Dimensions.get("window");
const STATUS_H =
  Platform.OS === "android" ? (StatusBar.currentHeight ?? 28) : 44;
const TOPBAR_H = 56;

// ─── Component ────────────────────────────────────────────────────────────────

export default function DriverMapScreen() {
  const [state, setState] = useState<ScreenState>("idle");
  const [activeTrip, setActiveTrip] = useState<Trip | null>(null);
  const cameraRef = useRef<Camera>(null);

  function flyTo(coord: [number, number], zoom = 14) {
    cameraRef.current?.setCamera({
      centerCoordinate: coord,
      zoomLevel: zoom,
      animationDuration: 600,
    });
  }

  function handleStartDelivery() {
    setState("pick_trip");
    flyTo([115.214, -8.667], 13);
  }

  function handleSelectTrip(trip: Trip) {
    setActiveTrip(trip);
    setState("active");
    flyTo(trip.coordinate, 15);
  }

  function handleCall() {
    Linking.openURL("tel:+628123456789");
  }

  function handleOpenMaps() {
    if (!activeTrip) return;
    const [lng, lat] = activeTrip.coordinate;
    Linking.openURL(
      `https://www.google.com/maps/dir/?api=1&destination=${lat},${lng}&travelmode=driving`,
    );
  }

  function handleFinish() {
    setActiveTrip(null);
    setState("pick_trip");
  }

  return (
    <View style={s.root}>
      <StatusBar barStyle="light-content" backgroundColor={BLUE} />

      {/* ── Top bar ── */}
      <View style={s.topBar}>
        <Text style={s.topName}>{DRIVER_NAME}</Text>
        <View style={s.topPkgRow}>
          <Text style={s.topPkgIcon}>⬡</Text>
          <Text style={s.topPkgCount}>{TOTAL_PACKAGES}</Text>
        </View>
      </View>

      {/* ── Mapbox MapView ── */}
      <MapView
        style={s.map}
        styleURL={Mapbox.StyleURL.Street}
        logoEnabled={false}
        attributionEnabled={false}
        compassEnabled={false}
      >
        <Camera
          ref={cameraRef}
          centerCoordinate={[115.214, -8.669]}
          zoomLevel={13}
          animationMode="flyTo"
        />

        {/* Polyline rute — tampil saat active */}
        {state === "active" && (
          <ShapeSource id="route" shape={ROUTE_GEOJSON}>
            <LineLayer
              id="routeLine"
              style={{
                lineColor: BLUE,
                lineWidth: 4,
                lineJoin: "round",
                lineCap: "round",
              }}
            />
          </ShapeSource>
        )}

        {/* Stop markers (circle + number) via MarkerView */}
        {(state === "pick_trip" || state === "active") &&
          TRIPS.map((trip, i) => (
            <MarkerView
              key={trip.id}
              coordinate={trip.coordinate}
              anchor={{ x: 0.5, y: 0.5 }}
            >
              <View
                style={[
                  s.markerDot,
                  activeTrip?.id === trip.id && s.markerDotActive,
                ]}
              >
                <Text
                  style={[
                    s.markerNum,
                    activeTrip?.id === trip.id && s.markerNumActive,
                  ]}
                >
                  {i + 1}
                </Text>
              </View>
            </MarkerView>
          ))}

        {/* Depot marker */}
        {state === "active" && (
          <MarkerView coordinate={DEPOT} anchor={{ x: 0.5, y: 0.5 }}>
            <View style={s.depotDot} />
          </MarkerView>
        )}
      </MapView>

      {/* ══ STATE: IDLE ══ */}
      {state === "idle" && (
        <View style={s.bottomIdle}>
          <TouchableOpacity
            style={s.startBtn}
            onPress={handleStartDelivery}
            activeOpacity={0.88}
            accessibilityRole="button"
            accessibilityLabel="Mulai pengantaran"
          >
            <View style={s.startBtnIcon}>
              <Text style={s.startBtnIconText}>▶</Text>
            </View>
            <Text style={s.startBtnLabel}>Mulai Pengantaran</Text>
          </TouchableOpacity>
        </View>
      )}

      {/* ══ STATE: PICK TRIP ══ */}
      {state === "pick_trip" && (
        <View style={s.bottomPickTrip}>
          <Text style={s.pickTripHint}>Pilih trip untuk dimulai</Text>
          <ScrollView
            horizontal
            showsHorizontalScrollIndicator={false}
            contentContainerStyle={s.tripScrollContent}
            decelerationRate="fast"
            snapToInterval={SW * 0.72 + 12}
            snapToAlignment="start"
          >
            {TRIPS.map((trip) => (
              <View key={trip.id} style={s.tripCard}>
                <View style={s.tripCardAccent} />
                <View style={s.tripCardBody}>
                  <Text style={s.tripCardTitle}>Trip {trip.number}</Text>
                  <Text style={s.tripCardRecipient}>{trip.recipientName}</Text>
                  <Text style={s.tripCardAddress}>{trip.address}</Text>
                  <Text style={s.tripCardMeta}>
                    {trip.distance} · {trip.duration}
                  </Text>
                  <TouchableOpacity
                    style={s.tripStartBtn}
                    onPress={() => handleSelectTrip(trip)}
                    activeOpacity={0.8}
                    accessibilityRole="button"
                    accessibilityLabel={`Mulai Trip ${trip.number}`}
                  >
                    <Text style={s.tripStartBtnText}>Mulai</Text>
                  </TouchableOpacity>
                </View>
              </View>
            ))}
          </ScrollView>
        </View>
      )}

      {/* ══ STATE: ACTIVE TRIP ══ */}
      {state === "active" && activeTrip && (
        <View style={s.bottomActive}>
          {/* Current trip card */}
          <View style={s.currentCard}>
            <View style={s.currentCardTop}>
              <View style={s.currentCardInfo}>
                <Text style={s.currentCardLabel}>Current Trip</Text>
                <View style={s.currentCardRow}>
                  <Text style={s.currentCardIcon}>👤</Text>
                  <Text style={s.currentCardRecipient}>
                    {activeTrip.recipientName}
                  </Text>
                </View>
                <View style={s.currentCardRow}>
                  <Text style={s.currentCardIcon}>📍</Text>
                  <Text style={s.currentCardAddress}>{activeTrip.address}</Text>
                </View>
                <Text style={s.currentCardPrice}>{activeTrip.price}</Text>
              </View>

              <View style={s.currentCardActions}>
                <TouchableOpacity
                  style={s.actionBtn}
                  onPress={handleCall}
                  activeOpacity={0.8}
                  accessibilityRole="button"
                  accessibilityLabel="Hubungi penerima"
                >
                  <Text style={s.actionBtnText}>📞 Call</Text>
                </TouchableOpacity>
                <TouchableOpacity
                  style={s.actionBtn}
                  activeOpacity={0.8}
                  accessibilityRole="button"
                  accessibilityLabel="Chat dengan penerima"
                >
                  <Text style={s.actionBtnText}>💬 Chat</Text>
                </TouchableOpacity>
              </View>
            </View>

            <TouchableOpacity
              style={s.distanceRow}
              onPress={handleOpenMaps}
              activeOpacity={0.8}
              accessibilityRole="button"
              accessibilityLabel={`Buka navigasi, jarak ${activeTrip.distance}`}
            >
              <Text style={s.distanceIcon}>▲</Text>
              <Text style={s.distanceText}>
                Jarak {activeTrip.distance} ({activeTrip.duration})
              </Text>
              <Text style={s.distanceChevron}>›</Text>
            </TouchableOpacity>
          </View>

          {/* Bottom sheet */}
          <View style={s.bottomSheet}>
            <View style={s.sheetHandle} />

            <View style={s.stepsArea}>
              {STEPS.map((step, i) => (
                <View key={i} style={s.stepRow}>
                  <View style={s.timelineCol}>
                    <View
                      style={[
                        s.stepDot,
                        step.done && s.stepDotDone,
                        step.active && s.stepDotActive,
                      ]}
                    />
                    {i < STEPS.length - 1 && (
                      <View style={[s.stepLine, step.done && s.stepLineDone]} />
                    )}
                  </View>
                  <Text
                    style={[
                      s.stepLabel,
                      step.active && s.stepLabelActive,
                      !step.done && !step.active && s.stepLabelFuture,
                    ]}
                  >
                    {step.label}
                  </Text>
                </View>
              ))}
            </View>

            <View style={s.sheetActions}>
              <TouchableOpacity
                style={s.photoBtn}
                activeOpacity={0.8}
                accessibilityRole="button"
                accessibilityLabel="Ambil foto bukti pengantaran"
              >
                <Text style={s.photoBtnText}>Foto Bukti</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={s.finishBtn}
                onPress={handleFinish}
                activeOpacity={0.8}
                accessibilityRole="button"
                accessibilityLabel="Selesaikan pengantaran"
              >
                <Text style={s.finishBtnText}>Selesai</Text>
              </TouchableOpacity>
            </View>
          </View>
        </View>
      )}
    </View>
  );
}

// ─── Styles ───────────────────────────────────────────────────────────────────

const s = StyleSheet.create({
  root: {
    flex: 1,
    backgroundColor: "#F4F4F4",
  },

  // Top bar
  topBar: {
    position: "absolute",
    top: 0,
    left: 0,
    right: 0,
    zIndex: 10,
    backgroundColor: BLUE,
    paddingTop: STATUS_H,
    height: STATUS_H + TOPBAR_H,
    paddingHorizontal: 20,
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
  },
  topName: {
    fontSize: 18,
    fontWeight: "700",
    color: WHITE,
    letterSpacing: -0.3,
  },
  topPkgRow: {
    flexDirection: "row",
    alignItems: "center",
    gap: 6,
  },
  topPkgIcon: {
    fontSize: 22,
    color: WHITE,
  },
  topPkgCount: {
    fontSize: 20,
    fontWeight: "700",
    color: WHITE,
  },

  // Map — mulai tepat di bawah top bar
  map: {
    flex: 1,
    marginTop: STATUS_H + TOPBAR_H,
  },

  // Markers
  markerDot: {
    width: 32,
    height: 32,
    borderRadius: 16,
    backgroundColor: BLUE_LIGHT,
    borderWidth: 2.5,
    borderColor: BLUE,
    alignItems: "center",
    justifyContent: "center",
  },
  markerDotActive: {
    backgroundColor: BLUE,
    borderColor: BLUE_DARK,
  },
  markerNum: {
    fontSize: 13,
    fontWeight: "700",
    color: BLUE_DARK,
  },
  markerNumActive: {
    color: WHITE,
  },
  depotDot: {
    width: 14,
    height: 14,
    borderRadius: 7,
    backgroundColor: WHITE,
    borderWidth: 3,
    borderColor: BLUE,
  },

  // ── IDLE ──
  bottomIdle: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
    paddingHorizontal: 20,
    paddingBottom: 36,
  },
  startBtn: {
    flexDirection: "row",
    alignItems: "center",
    backgroundColor: BLUE,
    borderRadius: 50,
    paddingVertical: 13,
    paddingLeft: 10,
    paddingRight: 28,
    elevation: 6,
  },
  startBtnIcon: {
    width: 44,
    height: 44,
    borderRadius: 22,
    backgroundColor: WHITE,
    alignItems: "center",
    justifyContent: "center",
    marginRight: 14,
  },
  startBtnIconText: {
    fontSize: 16,
    color: BLUE,
    marginLeft: 3,
  },
  startBtnLabel: {
    fontSize: 17,
    fontWeight: "600",
    color: WHITE,
    flex: 1,
    textAlign: "center",
    marginRight: 44,
  },

  // ── PICK TRIP ──
  bottomPickTrip: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
    paddingBottom: 28,
  },
  pickTripHint: {
    fontSize: 12,
    fontWeight: "500",
    color: TEXT_MID,
    paddingHorizontal: 20,
    marginBottom: 10,
  },
  tripScrollContent: {
    paddingHorizontal: 20,
    gap: 12,
    paddingBottom: 4,
  },
  tripCard: {
    width: SW * 0.72,
    backgroundColor: WHITE,
    borderRadius: 14,
    flexDirection: "row",
    overflow: "hidden",
    elevation: 4,
  },
  tripCardAccent: {
    width: 6,
    backgroundColor: BLUE,
  },
  tripCardBody: {
    flex: 1,
    padding: 14,
  },
  tripCardTitle: {
    fontSize: 12,
    fontWeight: "700",
    color: BLUE,
    textTransform: "uppercase",
    letterSpacing: 0.6,
    marginBottom: 4,
  },
  tripCardRecipient: {
    fontSize: 15,
    fontWeight: "600",
    color: TEXT_DARK,
    marginBottom: 2,
  },
  tripCardAddress: {
    fontSize: 13,
    color: TEXT_MID,
    lineHeight: 18,
    marginBottom: 6,
  },
  tripCardMeta: {
    fontSize: 12,
    color: TEXT_LIGHT,
    marginBottom: 10,
  },
  tripStartBtn: {
    backgroundColor: BLUE_LIGHT,
    borderRadius: 8,
    paddingVertical: 8,
    alignItems: "center",
  },
  tripStartBtnText: {
    fontSize: 14,
    fontWeight: "600",
    color: BLUE,
  },

  // ── ACTIVE ──
  bottomActive: {
    position: "absolute",
    bottom: 0,
    left: 0,
    right: 0,
  },
  currentCard: {
    backgroundColor: BLUE,
    marginHorizontal: 12,
    borderRadius: 16,
    overflow: "hidden",
    elevation: 6,
    marginBottom: 0,
  },
  currentCardTop: {
    flexDirection: "row",
    padding: 16,
    gap: 12,
  },
  currentCardInfo: {
    flex: 1,
  },
  currentCardLabel: {
    fontSize: 11,
    fontWeight: "700",
    color: "rgba(255,255,255,0.7)",
    textTransform: "uppercase",
    letterSpacing: 0.8,
    marginBottom: 6,
  },
  currentCardRow: {
    flexDirection: "row",
    alignItems: "flex-start",
    gap: 6,
    marginBottom: 3,
  },
  currentCardIcon: {
    fontSize: 13,
    marginTop: 1,
  },
  currentCardRecipient: {
    fontSize: 14,
    fontWeight: "500",
    color: WHITE,
    flex: 1,
  },
  currentCardAddress: {
    fontSize: 13,
    color: "rgba(255,255,255,0.85)",
    flex: 1,
    lineHeight: 18,
  },
  currentCardPrice: {
    fontSize: 22,
    fontWeight: "800",
    color: WHITE,
    marginTop: 8,
    letterSpacing: -0.5,
  },
  currentCardActions: {
    gap: 8,
    justifyContent: "flex-start",
    paddingTop: 18,
  },
  actionBtn: {
    backgroundColor: "rgba(255,255,255,0.18)",
    borderRadius: 10,
    paddingVertical: 8,
    paddingHorizontal: 14,
    alignItems: "center",
    minWidth: 88,
  },
  actionBtnText: {
    fontSize: 13,
    fontWeight: "600",
    color: WHITE,
  },
  distanceRow: {
    backgroundColor: WHITE,
    flexDirection: "row",
    alignItems: "center",
    paddingHorizontal: 16,
    paddingVertical: 12,
    gap: 8,
  },
  distanceIcon: {
    fontSize: 12,
    color: TEXT_MID,
  },
  distanceText: {
    fontSize: 13,
    fontWeight: "500",
    color: TEXT_DARK,
    flex: 1,
  },
  distanceChevron: {
    fontSize: 18,
    color: BLUE,
    fontWeight: "600",
  },

  // Bottom sheet
  bottomSheet: {
    backgroundColor: WHITE,
    borderTopLeftRadius: 20,
    borderTopRightRadius: 20,
    paddingTop: 10,
    paddingHorizontal: 20,
    paddingBottom: 32,
    elevation: 8,
    marginTop: 6,
  },
  sheetHandle: {
    width: 36,
    height: 4,
    backgroundColor: "#DDDDD8",
    borderRadius: 2,
    alignSelf: "center",
    marginBottom: 16,
  },
  stepsArea: {
    marginBottom: 18,
  },
  stepRow: {
    flexDirection: "row",
    alignItems: "flex-start",
    minHeight: 38,
  },
  timelineCol: {
    width: 24,
    alignItems: "center",
    marginRight: 12,
  },
  stepDot: {
    width: 12,
    height: 12,
    borderRadius: 6,
    backgroundColor: "#DDDDD8",
    marginTop: 3,
  },
  stepDotDone: {
    backgroundColor: BLUE,
  },
  stepDotActive: {
    backgroundColor: BLUE,
    borderWidth: 3,
    borderColor: BLUE_LIGHT,
    width: 14,
    height: 14,
    borderRadius: 7,
  },
  stepLine: {
    width: 2,
    flex: 1,
    backgroundColor: "#DDDDD8",
    minHeight: 24,
    marginVertical: 2,
  },
  stepLineDone: {
    backgroundColor: BLUE,
  },
  stepLabel: {
    fontSize: 14,
    fontWeight: "500",
    color: TEXT_DARK,
    paddingTop: 1,
    flex: 1,
    lineHeight: 20,
  },
  stepLabelActive: {
    fontWeight: "700",
  },
  stepLabelFuture: {
    color: TEXT_LIGHT,
    fontWeight: "400",
  },
  sheetActions: {
    flexDirection: "row",
    gap: 10,
  },
  photoBtn: {
    flex: 1,
    backgroundColor: BLUE_LIGHT,
    borderRadius: 10,
    paddingVertical: 13,
    alignItems: "center",
  },
  photoBtnText: {
    fontSize: 14,
    fontWeight: "600",
    color: BLUE,
  },
  finishBtn: {
    flex: 1,
    backgroundColor: BLUE_LIGHT,
    borderRadius: 10,
    paddingVertical: 13,
    alignItems: "center",
  },
  finishBtnText: {
    fontSize: 14,
    fontWeight: "600",
    color: BLUE,
  },
});
