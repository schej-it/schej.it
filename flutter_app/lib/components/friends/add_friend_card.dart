import 'package:flutter/material.dart';
import 'package:flutter_app/components/user_avatar.dart';
import 'package:flutter_app/constants/colors.dart';
import 'package:flutter_app/constants/constants.dart';
import 'package:flutter_app/constants/fonts.dart';
import 'package:flutter_app/models/api.dart';
import 'package:material_design_icons_flutter/material_design_icons_flutter.dart';
import 'package:provider/provider.dart';

class AddFriendCard extends StatelessWidget {
  final String id;
  final String name;
  final String picture;
  final String email;
  final bool requestAlreadySent;

  const AddFriendCard({
    Key? key,
    required this.id,
    required this.name,
    required this.picture,
    required this.email,
    this.requestAlreadySent = false,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: SchejConstants.listTileDecoration,
      child: Consumer<ApiService>(
        builder: (context, api, child) => ListTile(
          leading: UserAvatar(src: picture),
          title: Text(name, style: SchejFonts.subtitle),
          subtitle: Text(
            email,
            style: SchejFonts.body.copyWith(
              color: SchejColors.darkGray,
            ),
          ),
          trailing: IconButton(
            onPressed: () {
              api.sendFriendRequest(id);
            },
            icon: requestAlreadySent
                ? const Icon(MdiIcons.accountArrowRight,
                    color: SchejColors.gray)
                : const Icon(MdiIcons.accountPlus, color: SchejColors.black),
          ),
        ),
      ),
    );
  }
}
