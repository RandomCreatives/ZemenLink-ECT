import React, { useState, useEffect } from 'react';
import { StyleSheet, Text, View, FlatList, TouchableOpacity, SafeAreaView, TextInput, KeyboardAvoidingView, Platform } from 'react-native';
import { Send, Menu, Search, ShieldCheck, ArrowLeft } from 'lucide-react-native';
import { useMessaging } from './src/api/messaging';

const mockChats = [
  { id: '1', name: 'Enterprise Kernel Chat', lastMessage: 'Real-time sync enabled.', time: 'Just now', initials: 'ZL' },
  { id: '2', name: 'Compliance Team', lastMessage: 'Escrow request pending.', time: '12:45', initials: 'CT' },
];

export default function App() {
  const [token, setToken] = useState(null);
  const [activeChat, setActiveChat] = useState(null);
  const [inputText, setInputText] = useState('');
  const { messages, isConnected, sendMessage } = useMessaging(token);

  // Mock login for PoC
  useEffect(() => {
    const login = async () => {
      try {
        const response = await fetch('http://localhost:8080/auth/token', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ tenant_id: 'tenant-abc', user_id: 'mobile-user', role: 'user' })
        });
        const data = await response.json();
        setToken(data.token);
      } catch (err) {
        console.error("Auth failed", err);
      }
    };
    login();
  }, []);

  const handleSend = () => {
    if (inputText.trim()) {
      sendMessage(inputText);
      setInputText('');
    }
  };

  if (activeChat) {
    return (
      <SafeAreaView style={styles.container}>
        <KeyboardAvoidingView
          behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
          style={{ flex: 1 }}
        >
          {/* Chat Header */}
          <View style={styles.header}>
            <TouchableOpacity onClick={() => setActiveChat(null)}>
              <ArrowLeft color="#2481cc" size={24} />
            </TouchableOpacity>
            <Text style={styles.headerTitle}>{activeChat.name}</Text>
            <View style={{ width: 24 }} />
          </View>

          {/* Messages */}
          <FlatList
            data={messages}
            keyExtractor={(item) => item.id}
            renderItem={({ item }) => (
              <View style={[styles.messageBubble, item.sender_id === 'mobile-user' ? styles.selfMessage : styles.otherMessage]}>
                <Text style={styles.messageText}>{item.content}</Text>
              </View>
            )}
            contentContainerStyle={{ padding: 16 }}
          />

          {/* Input */}
          <View style={styles.inputContainer}>
            <TextInput
              style={styles.textInput}
              value={inputText}
              onChangeText={setInputText}
              placeholder="Write a message..."
            />
            <TouchableOpacity onPress={handleSend} disabled={!isConnected}>
              <Send color={isConnected ? "#2481cc" : "#ccc"} size={24} />
            </TouchableOpacity>
          </View>
        </KeyboardAvoidingView>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      {/* Header */}
      <View style={styles.header}>
        <Menu color="#2481cc" size={24} />
        <Text style={styles.headerTitle}>ZemenLink</Text>
        <Search color="#2481cc" size={24} />
      </View>

      {/* Chat List */}
      <FlatList
        data={mockChats}
        keyExtractor={(item) => item.id}
        renderItem={({ item }) => (
          <TouchableOpacity style={styles.chatItem} onPress={() => setActiveChat(item)}>
            <View style={styles.avatar}>
              <Text style={styles.avatarText}>{item.initials}</Text>
            </View>
            <View style={styles.chatInfo}>
              <View style={styles.chatHeader}>
                <Text style={styles.chatName}>{item.name}</Text>
                <Text style={styles.chatTime}>{item.time}</Text>
              </View>
              <Text style={styles.lastMessage} numberOfLines={1}>{item.lastMessage}</Text>
            </View>
          </TouchableOpacity>
        )}
      />

      {/* Bottom Status */}
      <View style={styles.footer}>
        <ShieldCheck color="#10b981" size={16} />
        <Text style={styles.footerText}>Siloed Isolation: tenant-abc</Text>
      </View>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
  },
  header: {
    height: 60,
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'space-between',
    paddingHorizontal: 16,
    borderBottomWidth: 1,
    borderBottomColor: '#eee',
  },
  headerTitle: {
    fontSize: 20,
    fontWeight: 'bold',
    color: '#2481cc',
  },
  chatItem: {
    flexDirection: 'row',
    padding: 16,
    borderBottomWidth: 0.5,
    borderBottomColor: '#eee',
    alignItems: 'center',
  },
  avatar: {
    width: 54,
    height: 54,
    borderRadius: 27,
    backgroundColor: '#2481cc',
    alignItems: 'center',
    justifyContent: 'center',
  },
  avatarText: {
    color: '#fff',
    fontSize: 20,
    fontWeight: 'bold',
  },
  chatInfo: {
    flex: 1,
    marginLeft: 16,
  },
  chatHeader: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    marginBottom: 4,
  },
  chatName: {
    fontSize: 16,
    fontWeight: 'bold',
  },
  chatTime: {
    fontSize: 12,
    color: '#999',
  },
  lastMessage: {
    fontSize: 14,
    color: '#666',
  },
  footer: {
    padding: 12,
    backgroundColor: '#f9f9f9',
    flexDirection: 'row',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 8,
  },
  footerText: {
    fontSize: 12,
    color: '#666',
  },
  messageBubble: {
    padding: 10,
    borderRadius: 15,
    marginBottom: 8,
    maxWidth: '80%',
  },
  selfMessage: {
    alignSelf: 'flex-end',
    backgroundColor: '#effdde',
  },
  otherMessage: {
    alignSelf: 'flex-start',
    backgroundColor: '#fff',
    borderWidth: 0.5,
    borderColor: '#eee',
  },
  messageText: {
    fontSize: 15,
  },
  inputContainer: {
    flexDirection: 'row',
    alignItems: 'center',
    padding: 12,
    borderTopWidth: 0.5,
    borderTopColor: '#eee',
    backgroundColor: '#fff',
  },
  textInput: {
    flex: 1,
    height: 40,
    backgroundColor: '#f1f1f1',
    borderRadius: 20,
    paddingHorizontal: 16,
    marginRight: 12,
  }
});
