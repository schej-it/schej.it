import 'dart:async';
import 'dart:collection';

import 'package:flutter/material.dart';
import 'package:flutter_app/components/app_bar.dart';
import 'package:flutter_app/components/friends/add_friend_card.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:flutter_app/models/user.dart';
import 'package:provider/provider.dart';

class AddFriendDialog extends StatefulWidget {
  final VoidCallback onClose;
  const AddFriendDialog({Key? key, required this.onClose}) : super(key: key);

  @override
  State<AddFriendDialog> createState() => _AddFriendDialogState();
}

class _AddFriendDialogState extends State<AddFriendDialog> {
  // Display variables.
  List<User> userSearchResults = [];

  // Delayed query variables.
  var callNum = 0;

  // Controllers.
  final TextEditingController _searchTextController = TextEditingController();

  // Dummy data.
  var results = [
    {'name': 'Winston Tilton', 'sent': false},
    {'name': 'Samantha Jones', 'sent': true},
    {'name': 'Tyler Smithson', 'sent': false},
    {'name': 'Arthi Singh', 'sent': false},
  ];

  @override
  void initState() {
    super.initState();

    // Start listening to changes.
    _searchTextController.addListener(delaySearch);
  }

  @override
  void dispose() {
    _searchTextController.dispose();
    super.dispose();
  }

  void delaySearch() {
    setState(() {
      callNum++;
    });
    final currCall = callNum;
    Timer(const Duration(milliseconds: 300), () {
      if (callNum == currCall) {
        _updateSearchResults();
      }
    });
  }

  Future<void> _updateSearchResults() async {
    if (_searchTextController.text == '') {
      setState(() {
        userSearchResults = [];
      });
    } else {
      ApiService api = context.read<ApiService>();
      // Replace all needed because semi-colon creates error in query, might want to move to backend later.
      final results = await api.getUserSearchResults(
          _searchTextController.text.replaceAll(';', '-'));
      setState(() {
        userSearchResults = results;
      });
    }
  }

  bool _isPotentialFriend(ApiService api, User u) {
    return api.authUser?.id != u.id && !api.friends.containsKey(u.id);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: SchejAppBar(
        isRoot: true,
        underline: false,
        centerTitle: true,
        titleString: 'Add friends',
        actions: [
          IconButton(onPressed: widget.onClose, icon: const Icon(Icons.close)),
        ],
      ),
      body: Container(
        color: SchejColors.white,
        padding: SchejConstants.pagePadding,
        child: Column(
          children: [
            _buildSearchTextField(),
            Expanded(
              child: Consumer<ApiService>(
                builder: (context, api, child) => _buildResults(api),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSearchTextField() {
    final textField = Padding(
      padding: const EdgeInsets.only(bottom: 20),
      child: TextField(
        autofocus: true,
        controller: _searchTextController,
        autocorrect: false,
        decoration: const InputDecoration(
          hintText: 'Search for a friend',
          prefixIcon: Icon(Icons.search),
        ),
        style: SchejFonts.subtitle,
      ),
    );

    return textField;
  }

  Widget _buildResults(ApiService api) {
    HashSet<String> outgoing = api.getOutgoingFriendRequestsUserIds();
    return ListView.builder(
      itemCount: userSearchResults.length,
      itemBuilder: (context, index) {
        final result = userSearchResults[index];
        return _isPotentialFriend(api, result)
            ? Padding(
                padding: const EdgeInsets.only(bottom: 10),
                child: AddFriendCard(
                  id: result.id,
                  name: result.fullName,
                  picture: result.picture,
                  email: result.email,
                  requestAlreadySent: outgoing.contains(result.id),
                ),
              )
            : const SizedBox(height: 0);
      },
    );
  }
}
